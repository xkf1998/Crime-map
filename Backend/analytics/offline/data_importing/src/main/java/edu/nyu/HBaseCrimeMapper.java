import ch.hsr.geohash.GeoHash;
import java.io.IOException;
import java.text.SimpleDateFormat;
import java.util.Date;
import org.apache.hadoop.hbase.client.Put;
import org.apache.hadoop.hbase.util.Bytes;
import org.apache.hadoop.io.LongWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Mapper;

public class HBaseCrimeMapper<K> extends Mapper<LongWritable, Text, K, Put> {
  static final byte[] EVENT_COLUMNFAMILY = Bytes.toBytes("e");
  static final byte[] LONGITUDE_QUALIFIER = Bytes.toBytes("x");
  static final byte[] LATITUDE_QUALIFIER = Bytes.toBytes("y");
  static final byte[] TIME_QUALIFIER = Bytes.toBytes("t");
  static final byte[] DESCRIPTION_QUALIFIER = Bytes.toBytes("d");

  private String normalize(double value, double minValue) {
    double posValue = value - minValue;
    long prefixValue = (long) posValue;
    String prefixStr = String.valueOf(prefixValue);
    // add extra '0' to the front
    String prefixPadding = new String(new char[3 - prefixStr.length()]).replace('\0', '0');
    prefixStr = prefixPadding + prefixStr;

    long suffixValue = (long) ((posValue - (double) prefixValue) * 1e6);
    String suffixStr = String.valueOf(suffixValue);
    // add extra '0' to the back
    String suffixPadding = new String(new char[6 - suffixStr.length()]).replace('\0', '0');
    suffixStr += suffixPadding;

    return prefixStr + suffixStr;
  }

  @Override
  public void map(LongWritable key, Text value, Context context)
      throws IOException, InterruptedException {
    String line = value.toString();
    if (line.length() == 0)
      return;

    String[] fields = line.split("\t");
    if (fields.length != 4) {
      context.getCounter("Bad Records", "BAD_FIELDS_NUM").increment(1);
      return;
    }

    // extract raw values of fields
    String rawLongitude = fields[0];
    String rawLatitude = fields[1];
    String timestamp = fields[2];
    String description = fields[3];

    String normalizedLongitude, normalizedLatitude, time, rowKey;

    // parse raw values and check whether they are valid
    SimpleDateFormat sdf = new SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss");
    try {
      // check whether values of latitude and longitude fields are valid double
      double longitude = Double.parseDouble(rawLongitude);
      double latitude = Double.parseDouble(rawLatitude);

      // It's impossible that a double value exactly equals to 0.0
      if (longitude == 0.0 || longitude < -180.0 || longitude > 180.0
          || latitude == 0.0 || latitude < -90.0 || latitude > 90.0) {
        context.getCounter("Bad Records", "MISSING_FIELD_VALUE").increment(1);
        return;
      }

      normalizedLongitude = normalize(longitude, -180.0);
      normalizedLatitude = normalize(latitude, -90.0);
      long unixTimestamp = Long.parseLong(timestamp) * 1000;
      Date date = new Date(unixTimestamp);
      time = sdf.format(date);
      GeoHash geohash = GeoHash.withCharacterPrecision(latitude, longitude, 12);
      rowKey = geohash.toBase32() + time;
    } catch (Exception e) {
      context.getCounter("Bad Records", "BAD_DATA_FORMAT:" + e).increment(1);
      return;
    }

    Put p = new Put(Bytes.toBytes(rowKey));

    p.addColumn(EVENT_COLUMNFAMILY, LONGITUDE_QUALIFIER, Bytes.toBytes(normalizedLongitude));
    p.addColumn(EVENT_COLUMNFAMILY, LATITUDE_QUALIFIER, Bytes.toBytes(normalizedLatitude));
    p.addColumn(EVENT_COLUMNFAMILY, TIME_QUALIFIER, Bytes.toBytes(time));
    p.addColumn(EVENT_COLUMNFAMILY, DESCRIPTION_QUALIFIER, Bytes.toBytes(description));
    context.write(null, p);
  }
}