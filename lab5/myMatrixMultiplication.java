import java.io.IOException;
import java.util.StringTokenizer;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;

public class MatrixMultiplication {

  public static class TokenizerMapperA
       extends Mapper<Object, Text, Text, String[]>{
    private Text word = new Text();

    public void map(Object key, Text value, Context context
                    ) throws IOException, InterruptedException {
      //StringTokenizer itr = new StringTokenizer(value.toString());
      String readLine = value.toString();
      String[] stringTokens = readLine.split(",");
      String idMatrix = stringTokens[0];
      String row = stringTokens[1];
      String col = stringTokens[2];
      String val = stringTokens[3];
      
      String[] values = {idMatrix, row, value};
      word.set(col);
      context.write(word, values);
    }
  }
  
  public static class TokenizerMapperB
       extends Mapper<Object, Text, Text, String[]>{
    private Text word = new Text();

    public void map(Object key, Text value, Context context
                    ) throws IOException, InterruptedException {
      //StringTokenizer itr = new StringTokenizer(value.toString());
      String readLine = value.toString();
      String[] stringTokens = readLine.split(",");
      String idMatrix = stringTokens[0];
      String row = stringTokens[1];
      String col = stringTokens[2];
      String val = stringTokens[3];
      
      String[] values = {idMatrix, col, value};
      word.set(row);
      context.write(word, values);
    }
  }
 
 
   public static class MultiplicationReducer
       extends Reducer<Text,String[],Text, IntWritable> {
    private IntWritable result = new IntWritable();

    public void reduce(Text key, Iterable<String[]> values,
                       Context context
                       ) throws IOException, InterruptedException {
      ArrayList<String[]> matrixA = new ArrayList<String[]>();
      ArrayList<String[]> matrixB = new ArrayList<String[]>();

      for(String[] element: values) {
      	String idMatrix = element[0];
      	//String row, col, val;
      	//val = element[2];
      	if idMatrix == "A"{
      		// then col = key
      		matrixA.add(element);
      	} else {
      		// then row = key
      		matrixB.add(element);
      	}
      }
      
      String rowA, colB;
      int valA, valB;
      String[] newKey; 
      
      for(String[] eleA: matrixA){
      	rowA = eleA[1];
      	valA = Integer.parseInt(eleA[2]);
      	for(String[] eleB: matrixB){
      		colB = eleB[1];
      		valB = Integer.parseInt(eleB[2]);
      		result.set(valA * valB);
      		newKey = {rowA, colB};
      		context.write(newKey, result);
      	}
      }

    }
  }


  public static class SecondMapper
       extends Mapper<Text, IntWritable, Text, IntWritable>{
    private Text word = new Text();

    public void map(Text key, IntWritable value, Context context
                    ) throws IOException, InterruptedException {
      context.write(key, value);
    }
  }  
  

  public static class SumReducer
       extends Reducer<Text,IntWritable,Text,IntWritable> {
    private IntWritable result = new IntWritable();

    public void reduce(Text key, Iterable<IntWritable> values,
                       Context context
                       ) throws IOException, InterruptedException {
                       
      int sum = 0;
      for (IntWritable val : values) {
        sum += val.get();
      }
      result.set(sum);
      context.write(key, result);
    }
  }

  public static void main(String[] args) throws Exception {
    Configuration conf = new Configuration();
    Job job = Job.getInstance(conf, "first map");
    job.setJarByClass(MatrixMultiplication.class);
    
    MultipleInputs.addInputPath(job, new Path(args[0]), TextInputFormat.class, TokenizerMapperA.class);
    MultipleInputs.addInputPath(job, new Path(args[1]), TextInputFormat.class, TokenizerMapperB.class);
    job.setReducerClass(MultiplicationReducer.class);
    job.setMapOutputValueClass(String[].class);
    job.setOutputKeyClass(Text.class);
    job.setReducerOutputValueClass(IntWritable.class);
    FileOutputFormat.setOutputPath(job, new Path(args[2]));
    job.waitForCompletion(true);
    
    Job job2 = Job.getInstance(conf, "second map");
    job2.setJarByClass(MatrixMultiplication.class);
    job2.setMapperClass(SecondMapper.class);
    job2.setReducerClass(SumReducer.class);
    job2.setMapOutputValueClass(String[].class);
    job2.setOutputKeyClass(Text.class);
    job2.setReducerOutputValueClass(IntWritable.class);
    
    FileInputFormat.setInputPaths(job2, new Path(args[2]));
    FileOutputFormat.setOutputPath(job2, new Path(args[2]));
    
    System.exit(job2.waitForCompletion(true) ? 0 : 1);
  }
}
