import org.apache.hadoop.conf.Configured;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.hbase.HBaseConfiguration;
import org.apache.hadoop.hbase.mapreduce.TableOutputFormat;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.util.Tool;
import org.apache.hadoop.util.ToolRunner;

/**
 * Uses HBase's {@link TableOutputFormat} to load crime data into a HBase table.
 */
public class HBaseCrimeImporter extends Configured implements Tool {
  @Override
  public int run(String[] args) throws Exception {
    if (args.length != 2) {
      System.err.println("Usage: HBaseCrimeImporter <input path> <table name>");
      return -1;
    }

    Job job = Job.getInstance();
    job.setJarByClass(HBaseCrimeImporter.class);
    job.setJobName("Data Importing");

    FileInputFormat.addInputPath(job, new Path(args[0]));
    job.getConfiguration().set(TableOutputFormat.OUTPUT_TABLE, args[1]);
    job.setMapperClass(HBaseCrimeMapper.class);

    job.setNumReduceTasks(0);
    job.setOutputFormatClass(TableOutputFormat.class);
    return job.waitForCompletion(true) ? 0 : 1;
  }

  public static void main(String[] args) throws Exception {
    int exitCode = ToolRunner.run(HBaseConfiguration.create(), new HBaseCrimeImporter(), args);
    System.exit(exitCode);
  }
}