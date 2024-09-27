
sample_0 = [[5, 1], [2, 1], [1, 1], [8, 1], [10, 0], [5, 0]]

def recurse(k, c):
    if not c:
        return 0
    if k == 0:
        return sum([v for v, important in c if not important]) + sum([-v for v, important in c if important])
    v, important = c[0]
    if not important:
        return v + recurse(k, c[1:])
    x = v + recurse(k - 1, c[1:])
    y = recurse(k, c[1:]) - v
    return max(x, y)

def luckBalance(k, c):
    c = sorted(c, key=lambda t: t[0])
    return recurse(k, c)

print(luckBalance(3, sample_0))

sample_02 = [[5, 1], [4, 0], [6, 1], [2, 1], [8, 0]]

print(luckBalance(2, sample_02))
