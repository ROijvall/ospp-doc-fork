import pygame
import resources as rs
import spritehandler

class Projectile(object):
    def __init__(self, projectileData):
        """ initialize function for the projectiles
        :param projectileData: interpreted Json data
        """
        self.x = projectileData["x"]
        self.y = projectileData["y"]
        self.a0 = projectileData["a"]


    def draw(self, win, projectileData):
        """ The drawing function for the projectiles
        :param win: the game window
        :param projectileData: interpreted Json data
        :return:
        """
        self.x = projectileData["x"]
        self.y = projectileData["y"]
        self.a = projectileData["a"]
        round(self.x)
        round(self.y)
        round(self.a)
        #oldRect = pipe.get_rect(center=(self.x, self.y))
        #rotatedBullet, newRect = rotate(bullet, oldRect, self.a)
        #win.blit(rotatedBullet, (self.x+33, self.y))
        pygame.draw.circle(win, (255,255,255), (round(self.x), round(self.y)), 4)



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


