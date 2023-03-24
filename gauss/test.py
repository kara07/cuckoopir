from scipy.stats import norm

x = ...  # 你想要计算概率的x值
cdf_value = norm.cdf(x, loc=0, scale=1)
probability = 1 - cdf_value
print(probability)
