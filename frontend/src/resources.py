import pygame, os, sys


"""
This module contains general constants and resources used in the program.
"""
pygame.init()
cwd = os.getcwd()       #\\src
frontFolder = os.path.dirname(cwd)    #\\frontend
print(frontFolder)
img = frontFolder + '\\img' #\\img

mapheight = 1200
mapwidth = 800

##SPRITES
tankSprite1 = [pygame.image.load(img + '\\LimeA.png'), pygame.image.load(img + '\\LimeB.png')]
tankSprite3 = [pygame.image.load(img + '\\OrangeA.png'), pygame.image.load(img + '\\OrangeB.png')]
tankSprite2 = [pygame.image.load(img + '\\GrapeA.png'), pygame.image.load(img + '\\GrapeB.png')]
tankSprite4 = [pygame.image.load(img + '\\BerryA.png'), pygame.image.load(img + '\\BerryB.png')]

explosion = pygame.image.load(img + '\\explosion-4.png')
jump = pygame.image.load(img + '\\preset_flat_fire.png')

pipe = pygame.image.load(img + '\\Pipe.png')
bg = pygame.image.load(img + '\\Snow2.png')
bg2 = pygame.image.load(img + '\\Woods.png')
bg = pygame.transform.scale(bg, (1200, 800))
bg2 = pygame.transform.scale(bg2, (1200, 800))
bullet = pygame.image.load(img + '\\Bullet.png')


##FONTS
largeText = pygame.font.Font('freesansbold.ttf',115)
smallText = pygame.font.Font('freesansbold.ttf',20)



##COLORS
black = (0,0,0)
white = (255,255,255)
grey = (125,125,125)
lightGrey = (190,190,190)
red = (200,0,0)
darkRed = (255,0,0)
green = (0,200,0)
darkGreen = (0,255,0)
lightPurple = (147,112,219)
purple = (138,43,226)
