import numpy as np
import time

N = 8192
A = np.random.rand(N, N)
B = np.random.rand(N, N)

start_time = time.time()
C = np.dot(A, B)
end_time = time.time()

elapsed_time = end_time - start_time
print(f"Matrix multiplication took {elapsed_time:.2f} seconds.")
