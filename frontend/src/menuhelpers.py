import pygame
import resources as rs

#tankCounter = 0
winX = 1200
winY = 700

tankCounter = 0
buttonColor = rs.darkRed
hooverColor = rs.red

def button(win, text, x, y, w, h, c, hc, action=None):
    """ Creates an interactive button
    :param win: game window
    :param text: Text shown on button
    :param x: x-position
    :param y: y-position
    :param w: Width
    :param h: Height
    :param c: Color
    :param hc: Hover color
    :param action: Pass along a certain action (like a function) for the button
    :return: Draws button graph
    """
    mouse = pygame.mouse.get_pos()
    click = pygame.mouse.get_pressed()
    if x+w > mouse[0] > x and y+h > mouse[1] > y:
        pygame.draw.rect(win,hc,(x,y,w,h))

        if click[0] == 1 and action != None:
            action()
    else:
        pygame.draw.rect(win,c,(x,y,w,h))

    TextSurf, TextRect = textObject(text, rs.smallText, rs.black)
    TextRect.center = ((x+(w/2)),(y+(h/2)))
    win.blit(TextSurf, TextRect)


def textObject(text, font, color):
    """ Function for rendering text
    :param text: Text that will be rendered
    :param font: Font of the text
    :return: Draws the main menu graphics
    """

    textSurface = font.render(text, True, color)
    return textSurface, textSurface.get_rect()


def chooseTank():
    global tankCounter
    global hooverColor
    global buttonColor
    if tankCounter == 0:
        tankCounter = 1
        buttonColor = rs.green
        hooverColor = rs.darkGreen
    elif tankCounter == 1:
        tankCounter = 2
        buttonColor = rs.lightPurple
        hooverColor = rs.purple
    elif tankCounter == 2:
        tankCounter = 3
        buttonColor = (255,127,80)
        hooverColor = (255,69,0)
    else:
        tankCounter = 0
        buttonColor = rs.red
        hooverColor = rs.darkRed







