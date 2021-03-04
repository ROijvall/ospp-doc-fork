import pygame
import json
import socket
import time
import errno

from tank import Tank
import projectile
import menuhelpers as mh
import resources as rs
import terrain
import spritehandler
import inputbox


pygame.init()
##################################
# constants needed for pygame
gamestate = ""
projectileDict = {}
tanksDict = {}
terrainDict = {}
explosionsDict = {}

BUFFER_SIZE = 131072
clock = pygame.time.Clock()
winX = 1200
winY = 700
win = pygame.display.set_mode((winX, winY))
pygame.display.set_caption("Super-Tanks")
tankCounter = 0
winX = 1200
winY = 700
health = 100
menu = True
choose = True
rs.bg.convert()
chosenBackground = rs.bg2

explosionEffect = spritehandler.Spritesheet(rs.explosion, 12, 1)
centerHandle = 4
index = 0

jumpFX = spritehandler.Spritesheet(rs.jump, 6, 5)
jumpIndex = 1
yo = 4

def mainMenu():
    """ The drawing function of the main menu
    :return: Draws the main menu graphics
    """
    global hooverColor
    global buttonColor
    global menu
    global chosenBackground
    #fade(1200, 700)

    menu = True
    global input_box1

    input_box1 = inputbox.InputBox((winX / 2) - 200, 500, 400, 32, "localhost")
    while menu:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                quit()
            input_box1.handle_event(event)
        win.blit(chosenBackground,(0, 0))
        TextSurf, TextRect = mh.textObject('Super-Tanks', rs.largeText, rs.red)
        TextRect.center = ((winX / 2), (winY / 2))
        win.blit(TextSurf, TextRect)

        mh.button(win, 'Play', (winX / 2) - 200, 550, 175, 50, rs.darkGreen, rs.green, chosenMap)
        mh.button(win, 'Quit', (winX / 2) + 25, 550, 175, 50, rs.darkRed, rs.red, quitGame)
        mh.button(win, 'Tank color', (winX / 2) - 200, 625, 400, 50, mh.buttonColor, mh.hooverColor, mh.chooseTank)
        #mh.button(win, 'Choose map', (winX / 2) + 50, 625, 200, 50, rs.grey, rs.lightGrey, chosenMap)
        pygame.draw.rect(win, (0, 0, 0), ((winX / 2) - 200, 500, 400, 32))
        input_box1.draw(win)

        pygame.display.update()


def quitGame():
    pygame.quit()
    quit()

def chosenMap():
    global menu
    global choose
    menu = False
    choose = True

    while choose:
        for event in pygame.event.get():
            if event.type == pygame.QUIT:
                pygame.quit()
                quit()
        win.blit(rs.bg, (0,0))
        win.blit(rs.bg2, (winX / 2, 0))

        TextSurf, TextRect = mh.textObject('Choose Map', rs.largeText, rs.red)
        TextRect.center = ((winX / 2), (winY / 2))
        win.blit(TextSurf, TextRect)

        mh.button(win, 'Winter', (winX / 2) - 250, 500, 175, 50, (135,206,250), (100,149,237), winterMap)
        mh.button(win, 'Forest', (winX / 2) + 75, 500, 175, 50, rs.darkGreen, (34,139,34), forestMap)

        pygame.display.update()

def winterMap():
    global chosenBackground
    chosenBackground = rs.bg
    gameLoop()

def forestMap():
    global chosenBackground
    chosenBackground = rs.bg2
    gameLoop()


def fade(width, height):
    fade = pygame.Surface((width, height))
    fade.fill((0, 0, 0))
    for alpha in range(0, 300):
        fade.set_alpha(alpha)
        win.fill((255,255,255))
        win.blit(fade, (0, 0))
        pygame.display.update()
        pygame.time.delay(1)

def jumping():
    global jumpIndex
    global yo
    if jumpIndex < 26:
        jumpFX.draw(win, index % jumpFX.totalCellCount, 200, 200)
        jumpIndex += 1
    else:
        jumpIndex = 0


def redrawGameWindow(gamestate):
    """ The main drawing function of the game
    :param gamestate: Decoded Json string
    :return: Draws the game graphics
    """
    global index
    global alive
    global uniqueID

    #print(gamestate["tanks"])
    if gamestate["tanks"][str(uniqueID)]["alive"] == False:
        alive = False

    win.blit(chosenBackground, (0, 0))
    projectiles = gamestate["projectile"]
    projList = []
    projToBeRemoved = []
    for key in projectiles:
        projList.append(key)
        if key in projectileDict:
            projectileDict[key].draw(win, projectiles[key])
        else:
            projectileDict[key] = projectile.Projectile(projectiles[key])
            projectileDict[key].draw(win, projectiles[key])

    for key in projectileDict:
        if key not in projList:
            projToBeRemoved.append(key)

    for key in projToBeRemoved:
        spriteCount = explosionEffect.draw(win, round(projectileDict[key].x), round(projectileDict[key].y),
                                           centerHandle)
        if spriteCount > 10:
            projectileDict.pop(key)


    tanks = gamestate["tanks"]
    for key in tanks:
        if key in tanksDict:
            tanksDict[key].draw(win, tanks[key])
        else:
            tanksDict[key] = Tank(tanks[key], mh.tankCounter)
            tanksDict[key].draw(win, tanks[key])
            if mh.tankCounter < 4:
                mh.tankCounter += 1
            else:
                mh.tankCounter = 0


    # generate terrain
    terrainMap = gamestate["terrain"]
    for key in terrainMap:
        terrainDict[key] = terrain.Terrain(terrainMap[key])
        if chosenBackground == rs.bg2:
            terrainDict[key].drawForest(win, terrainMap[key])
        else:
            terrainDict[key].drawWinter(win, terrainMap[key])


    pygame.display.update()

def connect(HOST, PORT, BUFFER, socket): #Rasmus - Kan lägga till IP här efter behov
    global uniqueID
    socket.connect((HOST, PORT))
    print("successful connect")
    jsonStr = socket.recv(BUFFER)
    #print(jsonStr)
    gamestate = json.loads(jsonStr.decode('utf-8'))
    uniqueID = gamestate["uniqueid"]
    print(uniqueID)
    redrawGameWindow(gamestate)
    socket.setblocking(0)
    return gamestate

def gameLoop():
    global choose
    global health
    global jumpIndex
    global input_box1
    global uniqueID
    global alive
    choose = False
    ticks_elapsed = 0
    #fade(1200,700)
    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        global gamestate
        gamestate = connect(input_box1.text, 8080, BUFFER_SIZE, s)
        ticks_elapsed = 0
        run = True
        alive = True
        """ main game loop
        """
        while run:
            clock.tick(27)
            ticks_elapsed += 1
            #if ticks_elapsed % 100:
            #    updateDict()

            for event in pygame.event.get():
                if event.type == pygame.QUIT:
                    run = False

            keys = pygame.key.get_pressed()
            str = ""
            if alive:
                if keys[pygame.K_d]: #drive left
                    str += "0,"
                elif keys[pygame.K_a]: #drive right
                    str += "1,"
                if keys[pygame.K_w]: #jump
                    str += "4,"
                if  keys[pygame.K_LEFT]: #rotate barrel left
                    str += "2,"
                elif keys[pygame.K_RIGHT]: #rotate barrel right
                    str += "3,"
                if keys[pygame.K_SPACE]: #shoot
                    str += "6,"
                elif keys[pygame.K_k] and health > 0: #currently useless
                    health -= 25
                if keys[pygame.K_m]:
                    print("hey i'm alive?")
                if keys[pygame.K_p]:
                    str += "9,"
                    print("sent death")
            if str != "":
                str = str.rstrip(',') # remove trailing comma
                str += "\n" # tells golang channels when the message is done
                s.sendall(str.encode())

            try:
                jsonData = s.recv(BUFFER_SIZE)
                try:
                    gamestate = json.loads(jsonData.decode('utf-8'))
                except ValueError as e:
                    print(e)
                redrawGameWindow(gamestate)
            except socket.error as e:
                if e.args[0] == errno.EWOULDBLOCK:
                    time.sleep(0)
                else:
                    print(e)
                    break

mainMenu()
pygame.quit()
