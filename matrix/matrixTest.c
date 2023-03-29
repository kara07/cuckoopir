#include <stdio.h>
#include <stddef.h>
#include <stdlib.h>
#include <stdint.h>
// #include <cblas.h>
#include <time.h>

typedef uint32_t Elem;
#define COMPRESSION 4
#define BASIS       10
#define BASIS2      BASIS*2
#define MASK        (1<<BASIS)-1

void print_matrix(const char *name, float *matrix, int rows, int cols){
    printf("%s:\n", name);
    for (int i = 0; i < rows; ++i){
        for (int j = 0; j < cols; ++j){
            printf("%f ", matrix[i * cols + j]);
        }
        printf("\n");
    }
}

void matMulVecPacked(Elem *out, const Elem *a, const Elem *b,
    size_t aRows, size_t aCols)
{
  Elem db, db2, db3, db4, db5, db6, db7, db8;
  Elem val, val2, val3, val4, val5, val6, val7, val8;
  Elem tmp, tmp2, tmp3, tmp4, tmp5, tmp6, tmp7, tmp8;
  size_t index = 0;
  size_t index2;

  for (size_t i = 0; i < aRows; i += 8) {
    tmp  = 0;
    tmp2 = 0;
    tmp3 = 0;
    tmp4 = 0;
    tmp5 = 0;
    tmp6 = 0;
    tmp7 = 0;
    tmp8 = 0;

    index2 = 0;
    for (size_t j = 0; j < aCols; j++) {
      db  = a[index];
      db2 = a[index+1*aCols];
      db3 = a[index+2*aCols];
      db4 = a[index+3*aCols];
      db5 = a[index+4*aCols];
      db6 = a[index+5*aCols];
      db7 = a[index+6*aCols];
      db8 = a[index+7*aCols];

      val  = db & MASK;
      val2 = db2 & MASK;
      val3 = db3 & MASK;
      val4 = db4 & MASK;
      val5 = db5 & MASK;
      val6 = db6 & MASK;
      val7 = db7 & MASK;
      val8 = db8 & MASK;
      tmp  += val*b[index2];
      tmp2 += val2*b[index2];
      tmp3 += val3*b[index2];
      tmp4 += val4*b[index2];
      tmp5 += val5*b[index2];
      tmp6 += val6*b[index2];
      tmp7 += val7*b[index2];
      tmp8 += val8*b[index2];
      index2 += 1;

      val  = (db >> BASIS) & MASK;
      val2 = (db2 >> BASIS) & MASK;
      val3 = (db3 >> BASIS) & MASK;
      val4 = (db4 >> BASIS) & MASK;
      val5 = (db5 >> BASIS) & MASK;
      val6 = (db6 >> BASIS) & MASK;
      val7 = (db7 >> BASIS) & MASK;
      val8 = (db8 >> BASIS) & MASK;
      tmp  += val*b[index2];
      tmp2 += val2*b[index2];
      tmp3 += val3*b[index2];
      tmp4 += val4*b[index2];
      tmp5 += val5*b[index2];
      tmp6 += val6*b[index2];
      tmp7 += val7*b[index2];
      tmp8 += val8*b[index2];
      index2 += 1;

      val  = (db >> BASIS2) & MASK;
      val2 = (db2 >> BASIS2) & MASK;
      val3 = (db3 >> BASIS2) & MASK;
      val4 = (db4 >> BASIS2) & MASK;
      val5 = (db5 >> BASIS2) & MASK;
      val6 = (db6 >> BASIS2) & MASK;
      val7 = (db7 >> BASIS2) & MASK;
      val8 = (db8 >> BASIS2) & MASK;
      tmp  += val*b[index2];
      tmp2 += val2*b[index2];
      tmp3 += val3*b[index2];
      tmp4 += val4*b[index2];
      tmp5 += val5*b[index2];
      tmp6 += val6*b[index2];
      tmp7 += val7*b[index2];
      tmp8 += val8*b[index2];
      index2 += 1;
      index += 1;
    }
    out[i]   += tmp;
    out[i+1] += tmp2;
    out[i+2] += tmp3;
    out[i+3] += tmp4;
    out[i+4] += tmp5;
    out[i+5] += tmp6;
    out[i+6] += tmp7;
    out[i+7] += tmp8;
    index += aCols*7;
  }
}

void matMulVec(Elem *out, const Elem *a, const Elem *b,
               size_t aRows, size_t aCols){
    Elem tmp;
    for (size_t i = 0; i < aRows; i++){
        tmp = 0;
        for (size_t j = 0; j < aCols; j++){
            tmp += a[aCols * i + j] * b[j];
        }
        out[i] = tmp;
    }
}

void matTransMul(Elem *out, const Elem *a, const Elem *b,
                 size_t aRows, size_t aCols, size_t bCols){
    for(size_t i = 0; i < aCols; i++){
        for(size_t j = 0; j < bCols; j++){
            for(size_t k = 0; k < aRows; k++){
                out[i * bCols + j] += a[k * aCols + i] * b[k * bCols + j];
            }
        }
    }
}

int main(){
    int rows = 10000;
    int cols = 10000;
    Elem *AI = malloc(rows * cols * sizeof(Elem));
    Elem *VI = malloc(cols * 1 * sizeof(Elem));
    Elem *out = malloc(rows * 1 * sizeof(Elem));

    srand(time(NULL));

    for (int i = 0; i < rows; i++){
        for (int j = 0; j < cols; j++){
            AI[i * cols + j] = (Elem)rand() / (Elem)(RAND_MAX);
        }
    }

    for (int i = 0; i < rows; i++){
        for (int j = 0; j < 1; j++){
            VI[i + j] = (Elem)rand() / (Elem)(RAND_MAX);
        }
    }

    srand(time(NULL));

    clock_t start = clock();
    matMulVecPacked(out, AI, VI, rows, cols);
    clock_t end = clock();
    double time_spent = (double)(end - start) / CLOCKS_PER_SEC;
    printf("Time spent for packed: %.4f seconds\n", time_spent);
    
    start = clock();
    matMulVec(out, AI, VI, rows, cols);
    end = clock();
    time_spent = (double)(end - start) / CLOCKS_PER_SEC;
    printf("Time spent for packed: %.4f seconds\n", time_spent);

    size_t rows_a = 3, cols_a = 2, cols_b = 2;

    Elem a[] = {1, 2, 3, 4, 5, 6}; // 3x2 矩阵
    Elem b[] = {1, 2, 3, 4, 5, 6}; // 3x2 矩阵

   Elem result[2 * 2] = {0};

    matTransMul(result, a, b, rows_a, cols_a, cols_b);

    printf("trans(a)xb：\n");
    for (int i = 0; i < cols_a; i++) {
        for (int j = 0; j < cols_b; j++) {
            printf("%u ", result[i * cols_b + j]);
        }
        printf("\n");
    }

    return 0;
}