# No machine should be able to reach another machine.
# This turns out to be a shortest-paths problem:
# we need to determine the shortest paths from a 
# given machine to every other, and eliminate the 
# the lowest-cost edge on each path. If a machine
# is connected to two other machines, let's say, 
# it won't suffice to eliminate an edge shared 
# by both the respective paths, since those two 
# machines would still be connected to each other.

class Graph:
    def __init__(self, cities, machines):
        self.adjacency = [None] * len(cities)
        self.cities = set(machines)

    def connect(self, u, v, cost):
        if self.adjacency[u] is None:
            self.adjacency[u] = {}
        if self.adjacency[v] is None:
            self.adjacency[v] = {}
        self.adjacency[u][v] = cost
        self.adjacency[v][u] = cost
        
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