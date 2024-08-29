# No machine should be able to reach another machine.
# This turns out to be a shortest-paths problem:
# we need to determine the shortest paths from a 
# given machine to every other, and eliminate the 
# the lowest-cost edge on each path. If a machine
# is connected to two other machines, let's say, 
# it won't suffice to eliminate an edge shared 
# by both the respective paths, since those two 
# machines would still be connected to each other.

import heapq
import os


class Graph:
    def __init__(self, cities):
        self.adjacency = {}
        for (u, v, cost) in cities:
            self.connect(u, v, cost)

    def connect(self, u, v, cost):
        if u not in self.adjacency:
            self.adjacency[u] = {}
        if v not in self.adjacency:
            self.adjacency[v] = {}
        self.adjacency[u][v] = cost
        self.adjacency[v][u] = cost
        
    def disconnect(self, u, v):
        del(self.adjacency[v], u)
        del(self.adjacency[u], v)
        
    def find_shortest_path(self, source, target):
        queue = []
        distances = {}
        previous = {}
        visited = set()
        
        distances[source] = 0
        heapq.heappush(queue, (0, source))

        while queue:
            _, u = heapq.heappop(queue)
            if u == target:
                break
            for v in self.adjacency[u]:
                if v in visited:
                    continue
                alt = distances[u] + self.adjacency[u][v]
                if v not in distances:
                    distances[v] = alt
                    previous[v] = u
                elif distances[v] > alt:
                    distances[v] = alt
                    previous[v] = u
                heapq.heappush(queue, (distances[v], v))
                
            visited.add(u)

        v = target
        path = []
        while v != source:
            path = [v] + path
            v = previous[v]
        return [source] + path
    
def minTime(roads, machines) -> int:
    graph = Graph(roads)
    cheapest_edges = []
    for machine in machines:
        for other_machine in [m for m in machines if m > machine]:
            previous = graph.find_shortest_path(machine, other_machine)
            # What is the lowest-cost edge on this path?
            edges = [(previous[i - 1], previous[i]) for i in range(1, len(previous))]
            cheapest = min(edges, key = lambda e: graph.adjacency[e[0]][e[1]])
            if cheapest[0] > cheapest[1]:
                cheapest = (cheapest[1], cheapest[0])
            cheapest_edges.append(cheapest)
    s = set(cheapest_edges)
    costs = [graph.adjacency[u][v] for u, v in s]
    result = sum(costs)
    return result


# Sample 0
result0 = minTime([[2, 1, 8], [1, 0, 5], [2, 4, 5], [1, 3, 4]], [2, 4, 0])
# Sample 1
result1 = minTime([[0, 1, 4], [1, 2, 3], [1, 3, 7], [0, 4, 2]], [2, 3, 4])
# Sample 2
result2 = minTime([[0, 3, 3], [1, 4, 4], [1, 3, 4], [0, 2, 5]], [1, 3, 4])

print(result0, result1, result2)
