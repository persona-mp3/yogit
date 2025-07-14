public class BinarySearch {
  public static void main(String[] args) {
    int[] nums = {2, 3, 20, 82, 124, 300 };
    int start = 0;
    int end = nums.length - 1;
    int target = 11;
    BinaryHelper(nums, start, end, target);
  }

  static int BinaryHelper(int[] data, int start, int end, int target) {
    int middle = (start+end)/2;
    
    if (start > end) {
      System.out.printf("target not found, start is more than end\n");
      return -1;
    }

    if (data[middle] == target) {
      System.out.printf("target found at position %d with element %d\n", middle, data[middle]);
      return middle;

    } else if (data[middle] > target) {
      BinaryHelper(data, start, (middle - 1), target);
    } else {
      BinaryHelper(data, (middle + 1), end, target);
    }

    System.out.printf("target not found\n");
    return -1;
  }
}
