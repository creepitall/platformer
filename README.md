# platformer

simple 2D platformer with golang/pixel(OpenGL 3.3)

master **ISN'T** actual branch

---
Current video demo:
https://youtu.be/OGulTJXcpyg

---

12/01/22
- global fix architecture(add object Character, add scene, delete global param etc.) 
- add load config 
- add save log file for test

28/11/21
- fix character camera (in height)
- fix character movement outside the screen 
- add screen metrics (rectangle highlights for physical object)
- add physical objects

27/11/21
- fix character camera
- add screen metrics
- add new front layer

26/11/21
- animation run, idle, jump (jump is have some lags)
- fix FPS (limit with 60 FPS)
- reduced GPU/CPU load
- change architecture app, load assets scene and character became more fast
- add character status (dead, jump)
- add trigger of death, if status true - character reset to start point
- delete double/triple/etc. jump
- add reset for a key

15/11/12
- add main windows
- add load assets and create sprites
- add grab scene and camera
- add some layer (background, front, character)
