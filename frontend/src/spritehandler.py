import pygame
import resources

class Spritesheet(object):
    def __init__(self, filename, cols, rows):
        self.sheet = filename

        self.cols = cols
        self.rows = rows
        self.totalCellCount = cols * rows
        self.cellIndex = 0

        self.rect = self.sheet.get_rect()
        w = self.cellWidth = self.rect.width / cols
        h = self.cellHeight = self.rect.height / rows
        hw, hh = self.cellCenter = (w / 2, h / 2)

        self.cells = list([(index % cols * w, index / cols * h, w, h) for index in range(self.totalCellCount)])
        self.handle = list([
            (0, 0), (-hw, 0), (-w, 0),
            (0, -hh), (-hw, -hh), (-w, -hh),
            (0, -h), (-hw, -h), (-w, -h), ])

    def draw(self, surface, x, y, handle=0):
        index = self.cellIndex % self.totalCellCount
        if index <= self.totalCellCount:
            surface.blit(self.sheet, (x + self.handle[handle][0], y + self.handle[handle][1]), self.cells[index])
            self.cellIndex += 1
        return index