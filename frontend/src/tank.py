import pygame
import resources as rs

health = 175
tanksDict = ""

class Tank(object):
    def __init__(self, tankData, colorIndex):
        """ Initialize function for the tanks
        :param tankData: Interpreted Json data
        """
        self.x = tankData["x"]
        self.y = tankData["y"]
        self.width = 32 #  For hitbox, random number
        self.height = 32 # For hitbox, random number
        self.degTank = tankData["degTank"]
        self.degCannon = tankData["degCannon"]
        self.right = False
        self.left = False
        self.xCannon = tankData["x"]
        self.yCannon = tankData["y"]
        self.isMoving = True
        self.walkCount = 0
        self.dir = tankData["dir"]
        self.hp = tankData["hp"]
        self.alive = tankData["alive"]
        self.colorIndex = colorIndex

    def draw(self, win, tankData):
        """ Draws the tank on the window surface pointing in the right direction.
        Chooses the right sprite depending on the amount of "ticks" with the
        help of walkCount.
        Also handles the angle of the pipe and tank according to the surface and
        aim of the player.
        :param win: the game window
        :param tankData: The interpreted Json data
        :return: draws the tank on the window surface
        """
        self.x = tankData["x"]
        self.y = tankData["y"] - 10
        self.degCannon = tankData["degCannon"]
        self.degTank = tankData["degTank"]
        self.dir = tankData["dir"]
        self.hp = tankData["hp"]
        self.alive = tankData["alive"]

        if self.alive:
            if self.isMoving is True and self.walkCount < 10:
                self.walkCount += 1
            else:
                self.walkCount = 0

            if self.colorIndex == 1:
                colorSprite = rs.tankSprite1
            elif self.colorIndex == 2:
                colorSprite = rs.tankSprite2
            elif self.colorIndex == 3:
                colorSprite = rs.tankSprite3
            else:
                colorSprite = rs.tankSprite4


            self.xCannon = self.x
            self.yCannon = self.y - 2
            oldRect = rs.pipe.get_rect(center=(self.xCannon, self.yCannon))
            rotatedPipe, newRect = rotate(rs.pipe, oldRect, self.degCannon-self.degTank)
            win.blit(rotatedPipe, newRect)



            if self.walkCount < 5:
                oldRectTtank = colorSprite[0].get_rect(center=(self.x, self.y - 8))
                rotatedTank, newTankRect = rotate(colorSprite[0], oldRectTtank, self.degTank*-1)
                win.blit(rotatedTank, newTankRect)

            else:
                oldRectTtank = colorSprite[1].get_rect(center=(self.x, self.y - 8))
                rotatedTank, newTankRect = rotate(colorSprite[1], oldRectTtank, self.degTank*-1)
                win.blit(rotatedTank, newTankRect)

            pygame.draw.rect(win, (255,0,0), (self.x - 25, self.y - 50, 50, 5))
            pygame.draw.rect(win, (0, 250, 0), (self.x - 25, self.y - 50, self.hp//2, 5))

def rotate(image, rect, angle):
    """ rotates an image while maintaining its position
    :param image: A loaded image
    :param rect: A rect surface of the original image
    :param angle: The angle of the rotation
    :return: rotated image,
    """
    rot_image = pygame.transform.rotate(image, angle)
    rot_rect = rot_image.get_rect(center=rect.center)
    return rot_image, rot_rect
