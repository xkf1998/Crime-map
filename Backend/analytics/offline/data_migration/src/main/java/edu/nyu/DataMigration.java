import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.NullWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class DataMigration {

  public static void main(String[] args) throws Exception {
    if (args.length < 2) {
      System.err.println("Usage: DataMigration <output path> <input path 1> <input path 2> ... <input path n>");
      System.exit(-1);
    }

    Job job = Job.getInstance();
    job.setJarByClass(DataMigration.class);
    job.setJobName("Data Cleaning");

    FileOutputFormat.setOutputPath(job, new Path(args[0]));
    for (int i = 1; i < args.length; i++) {

    }

    job.setReducerClass(DataMigrationReducer.class);

    job.setNumReduceTasks(1);
    job.setOutputKeyClass(Text.class);
    job.setOutputValueClass(NullWritable.class);
    job.setMapOutputKeyClass(NullWritable.class);
    job.setMapOutputValueClass(Text.class);

    System.exit(job.waitForCompletion(true) ? 0 : 1);
  }
}
