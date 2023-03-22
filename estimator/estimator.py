from estimator import *;
import math
# n is the lattice dimension that is used
for n in [1024]:
    print(LWE.primal_usvp (
        LWE.Parameters (
        n = n,
        q = 2**32,
        Xs = ND.Uniform(-1 ,1),
        # Xe = ND.Uniform(-1 ,1),
        Xe = ND.DiscreteGaussian(6.4),
        
        m = 1024
    )
))