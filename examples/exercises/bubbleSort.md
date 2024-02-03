### Bubble Sort Algorithm

<br />

Bubble sort is a basic method that sorts a list by comparing each element and swapping them if needed. This code has a function called `bubbleSort`. It organizes a list of numbers from smallest to largest.

<br />

Here's how the code works:

- The function `bubbleSort` starts by setting up variables **i** and **j** for looping, and **arrSize** to store the length of the array.

- It then enters a loop running from the start of the array to two less than its length `(n - 2)`. This adjustment is made because the comparison involves looking ahead by one element `(arr[j + 1])`, and we want to avoid going out of bounds.

- Inside this loop, another loop starts, running from `0` up to `n - i - 2`. The `- i` part ensures that with each outer loop iteration, one fewer element is compared since the largest elements gradually bubble up to the end of the array and do not need to be compared again.

- Within the inner loop, adjacent elements are compared. If an element at position `j` is greater than the one next to it `(j + 1)`, they are swapped. This is done using a temporary variable temp to hold one of the values during the swap.

- The function then returns the sorted array.
