#include "matrix_multiply.h"
#include <stdio.h>
// #include <cblas.h>

void matMul(Elem *out, const Elem *a, const Elem *b,
            size_t aRows, size_t aCols, size_t bCols){
    for (size_t i = 0; i < aRows; i++){
        for (size_t k = 0; k < aCols; k++){
            for (size_t j = 0; j < bCols; j++){
                out[bCols * i + j] += a[aCols * i + k] * b[bCols * k + j];
            }
        }
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
            for(size_t k = 0; k < aRows; k++){
                for(size_t j = 0; j < bCols; j++){
                    out[i * bCols + j] += a[k * aCols + i] * b[k * bCols + j];
            }
        }
    }
}

void transpose(Elem *out, const Elem *in, size_t rows, size_t cols){
    for (size_t i = 0; i < rows; i++){
        for (size_t j = 0; j < cols; j++){
            out[j * rows + i] = in[i * cols + j];
        }
    }
}
