import pygame
import resources as rs

class Terrain(object):

    def __init__(self, terrainData):
        self.x = terrainData["x"]
        self.y = terrainData["y"]
        #self.material = terrainData["material"]
        #self.canFall = terrainData["canfall"]

    def drawForest(self, win, terrainData):
        self.x = terrainData["x"]
        self.y = terrainData["y"]

        pygame.draw.rect(win, (40,75,0), (self.x, self.y, 1, 1))
        pygame.draw.rect(win, (40,75,0), (self.x, self.y+1, 1, rs.mapwidth-self.y))

    def drawWinter(self, win, terrainData):
        self.x = terrainData["x"]
        self.y = terrainData["y"]

        pygame.draw.rect(win, (135,206,235), (self.x, self.y, 1, 1))
        pygame.draw.rect(win, (135,206,235), (self.x, self.y+1, 1, rs.mapwidth-self.y))

