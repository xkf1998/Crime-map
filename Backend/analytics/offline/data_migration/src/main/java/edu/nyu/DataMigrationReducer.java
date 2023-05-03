import java.io.IOException;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Reducer;

public class DataMigrationReducer
    extends Reducer<NullWritable, Text, Text, NullWritable> {

  @Override
  public void reduce(NullWritable key, Iterable<Text> records, Context context)
      throws IOException, InterruptedException {
    for (Text record : records) {
      context.write(record, NullWritable.get());
    }
  }
}
