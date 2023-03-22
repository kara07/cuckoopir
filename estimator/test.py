load("./estimator.py")
n, alpha, q = 256, 0.000976562500000000, 65537
set_verbose(1)
_ = estimate_lwe(n, alpha, q)
