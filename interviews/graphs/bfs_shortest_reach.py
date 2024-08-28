# The problem is simply to find all the shortest paths from a particular node
# to every other...i.e., to implement Dijkstra's algorithm.

class Graph:
    def __init__(self, order):
        self.adjacency = [None] * order

    def connect(self, u, v):
        if self.adjacency[u] is None:
            self.adjacency[u] = set()
        if self.adjacency[v] is None:
            self.adjacency[v] = set()
        self.adjacency[u].add(v)
        self.adjacency[v].add(u)
        
    def find_all_distances(self, source):
        queue = set()
        distances = {}
        visited = set()
        queue.add(source)
        while queue:
            u = queue.pop(0)
            for v in self.adjacency[v]:
                if v in visited:
                    continue
                queue.add(v)
                alt = distances[u] + 1
                if v not in distances:
                    distances[v] = alt
                else:
                    distances[v] += 1
                
            visited.add(u)
            
        return [distances.get(v, -1) for v in range(len(self.adjacency))]


t = int(input())
for i in range(t):
    n,m = [int(value) for value in input().split()]
    graph = Graph(n)
    for i in range(m):
        x,y = [int(x) for x in input().split()]
        graph.connect(x-1,y-1) 
    s = int(input())
    graph.find_all_distances(s-1)