import java.io.IOException;
import java.util.StringTokenizer;
import java.util.ArrayList;
import java.lang.Object;

import java.io.DataInput;
import java.io.DataOutput;

import org.apache.hadoop.conf.Configuration;
import org.apache.hadoop.fs.Path;
import org.apache.hadoop.io.IntWritable;
import org.apache.hadoop.io.Writable;
import org.apache.hadoop.io.WritableComparable;
import org.apache.hadoop.io.Text;
import org.apache.hadoop.mapreduce.Job;
import org.apache.hadoop.mapreduce.Mapper;
import org.apache.hadoop.mapreduce.Reducer;
import org.apache.hadoop.mapreduce.lib.input.FileInputFormat;
import org.apache.hadoop.mapreduce.lib.output.FileOutputFormat;
import org.apache.hadoop.mapreduce.lib.input.TextInputFormat;
import org.apache.hadoop.mapreduce.lib.input.MultipleInputs;
import org.apache.hadoop.util.ReflectionUtils;


class CellWritable implements Writable {
	char idMatrix;
	int position;
	int value;
	
	CellWritable(){
		this.idMatrix = 'A';
		this.position = 0;
		this.value = 0;
	}
	
	CellWritable(char id, int p, int v){
		this.idMatrix = id;
		this.position = p;
		this.value = v;
	}
	
	@Override
	public void readFields(DataInput input) throws IOException {
		idMatrix = input.readChar();
		position = input.readInt();
		value = input.readInt();
	}
	
	@Override
	public void write(DataOutput output) throws IOException {
		output.writeChar(idMatrix);
		output.writeInt(position);
		output.writeInt(value);
	}
	
	
}


class PairWritable implements WritableComparable<PairWritable> {
	int row;
	int col;
	
	PairWritable(){
		this.row = 0;
		this.col = 0;
	}
	
	PairWritable(int row, int col){
		this.row = row;
		this.col = col;
	}
	
	@Override
	public void readFields(DataInput input) throws IOException {
		row = input.readInt();
		col = input.readInt();
	}
	
	@Override
	public void write(DataOutput output) throws IOException {
		output.writeInt(row);
		output.writeInt(col);
	}
	
	public String toString() {
		return this.row + " " + this.col + " ";
	}
	
	@Override
	public int compareTo(PairWritable compare) {
		
		if (this.row > compare.row) {
			return 1;
		} else if (this.row < compare.row) {
			return -1;
		} else {
			if(this.col > compare.col) {
				return 1;
			} else if (this.col < compare.col) {
				return -1;
			}
		}
		return 0;
	}
}


public class MatrixMultiplication {

  public static class TokenizerMapperA
       extends Mapper<Object, Text, Text, CellWritable>{
    private Text word = new Text();

    public void map(Object key, Text value, Context context
                    ) throws IOException, InterruptedException {

      String readLine = value.toString();
      String[] stringTokens = readLine.split(" ");
      char idMatrix = stringTokens[0].charAt(0);
      int row = Integer.parseInt(stringTokens[1]);
      String col = stringTokens[2];
      int val = Integer.parseInt(stringTokens[3]);
      
      CellWritable values = new CellWritable(idMatrix, row, val);
      word.set(col);
      context.write(word, values);
    }
  }
  
  public static class TokenizerMapperB
       extends Mapper<Object, Text, Text, CellWritable>{
    private Text word = new Text();

    public void map(Object key, Text value, Context context
                    ) throws IOException, InterruptedException {
      //StringTokenizer itr = new StringTokenizer(value.toString());
      String readLine = value.toString();
      String[] stringTokens = readLine.split(" ");
      char idMatrix = stringTokens[0].charAt(0);
      String row = stringTokens[1];
      int col = Integer.parseInt(stringTokens[2]);
      int val = Integer.parseInt(stringTokens[3]);
      
      CellWritable values = new CellWritable(idMatrix, col, val);
      word.set(row);
      context.write(word, values);
    }
  }
 
 
   public static class MultiplicationReducer
       extends Reducer<Text, CellWritable, PairWritable, IntWritable> {
    private IntWritable result = new IntWritable();

    public void reduce(Text key, Iterable<CellWritable> values,
                       Context context
                       ) throws IOException, InterruptedException {
      ArrayList<CellWritable> matrixA = new ArrayList<CellWritable>();
      ArrayList<CellWritable> matrixB = new ArrayList<CellWritable>();
      
      Configuration conf = context.getConfiguration();

      for(CellWritable elem: values) {
      
      	// La gestió de memòria a Java és una locura
      	CellWritable cell = ReflectionUtils.newInstance(CellWritable.class, conf);
      	ReflectionUtils.copy(conf, elem, cell);
      	
      	if (cell.idMatrix == 'A'){
      		matrixA.add(cell);
      	} else if (cell.idMatrix == 'B'){
      		matrixB.add(cell);
      	}
      }
      
      for(CellWritable eleA: matrixA){
      	for(CellWritable eleB: matrixB){
      		int mult = eleA.value * eleB.value;
      		result.set(mult);
      		PairWritable newKey = new PairWritable (eleA.position, eleB.position);
      		context.write(newKey, result);
      	}
      }

    }
  }


  public static class SecondMapper
       extends Mapper<Object, Text, PairWritable, IntWritable>{

    public void map(Object key, Text value, Context context
                    ) throws IOException, InterruptedException {
      String readLine = value.toString();
      String[] stringTokens = readLine.split(" ");

      int row = Integer.parseInt(stringTokens[0]);
      int col = Integer.parseInt(stringTokens[1]);
      int valint = Integer.parseInt(stringTokens[2].replaceAll("\\s+",""));
      PairWritable pair = new PairWritable(row, col);
      IntWritable val = new IntWritable(valint);
      
      context.write(pair, val);
    }
  }  
  

  public static class SumReducer
       extends Reducer<PairWritable,IntWritable,PairWritable,IntWritable> {
    private IntWritable result = new IntWritable();

    public void reduce(PairWritable key, Iterable<IntWritable> values,
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
    job.setMapOutputKeyClass(Text.class);
    job.setOutputKeyClass(PairWritable.class);
    job.setMapOutputValueClass(CellWritable.class);


    job.setOutputValueClass(IntWritable.class);
    FileOutputFormat.setOutputPath(job, new Path(args[2]));
    job.waitForCompletion(true);
    
    Job job2 = Job.getInstance(conf, "second map");
    job2.setJarByClass(MatrixMultiplication.class);
    job2.setMapperClass(SecondMapper.class);
    job2.setReducerClass(SumReducer.class);
    job2.setOutputValueClass(IntWritable.class);
    job2.setOutputKeyClass(PairWritable.class);
    
    FileInputFormat.setInputPaths(job2, new Path(args[2]));
    FileOutputFormat.setOutputPath(job2, new Path(args[3]));
    
    System.exit(job2.waitForCompletion(true) ? 0 : 1);
  }
}
