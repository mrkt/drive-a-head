// 虚拟摇杆控制
(function() {
  var isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
  if (!isMobile) return;

  // 等待 Unity canvas 加载
  function initJoystick() {
    var canvas = document.querySelector("#unity-canvas");
    if (!canvas) {
      setTimeout(initJoystick, 100);
      return;
    }

    var joystick = document.getElementById('virtual-joystick');
    var joystickStick = document.getElementById('joystick-stick');
    var joystickBase = document.getElementById('joystick-base');

    if (!joystick || !joystickStick || !joystickBase) {
      console.error('Virtual joystick elements not found');
      return;
    }

    joystick.style.display = 'block';

    var currentKeys = {
      horizontal: 0, // -1 左, 0 中, 1 右
      vertical: 0    // -1 下, 0 中, 1 上
    };

    var keyStates = {
      'ArrowUp': false,
      'ArrowDown': false,
      'ArrowLeft': false,
      'ArrowRight': false,
      'w': false,
      's': false,
      'a': false,
      'd': false
    };

    function updateDirectionIndicators(horizontal, vertical) {
      var indicators = document.querySelectorAll('.direction-indicator');
      indicators.forEach(function(ind) { ind.classList.remove('active'); });

      if (vertical > 0 && horizontal === 0) {
        document.querySelector('.dir-up').classList.add('active');
      } else if (vertical < 0 && horizontal === 0) {
        document.querySelector('.dir-down').classList.add('active');
      } else if (horizontal < 0 && vertical === 0) {
        document.querySelector('.dir-left').classList.add('active');
      } else if (horizontal > 0 && vertical === 0) {
        document.querySelector('.dir-right').classList.add('active');
      } else if (vertical > 0 && horizontal < 0) {
        document.querySelector('.dir-up-left').classList.add('active');
      } else if (vertical > 0 && horizontal > 0) {
        document.querySelector('.dir-up-right').classList.add('active');
      } else if (vertical < 0 && horizontal < 0) {
        document.querySelector('.dir-down-left').classList.add('active');
      } else if (vertical < 0 && horizontal > 0) {
        document.querySelector('.dir-down-right').classList.add('active');
      }
    }

    function simulateKeyEvent(key, isDown) {
      var event = new KeyboardEvent(isDown ? 'keydown' : 'keyup', {
        key: key,
        code: key === 'ArrowUp' ? 'ArrowUp' :
              key === 'ArrowDown' ? 'ArrowDown' :
              key === 'ArrowLeft' ? 'ArrowLeft' :
              key === 'ArrowRight' ? 'ArrowRight' :
              'Key' + key.toUpperCase(),
        keyCode: key === 'ArrowUp' ? 38 :
                 key === 'ArrowDown' ? 40 :
                 key === 'ArrowLeft' ? 37 :
                 key === 'ArrowRight' ? 39 :
                 key.toUpperCase().charCodeAt(0),
        bubbles: true,
        cancelable: true
      });
      canvas.dispatchEvent(event);
    }

    function updateKeys(horizontal, vertical) {
      // 释放所有之前的按键
      for (var key in keyStates) {
        if (keyStates[key]) {
          simulateKeyEvent(key, false);
          keyStates[key] = false;
        }
      }

      // 按下新的按键
      if (vertical > 0) { // 上
        simulateKeyEvent('ArrowUp', true);
        simulateKeyEvent('w', true);
        keyStates['ArrowUp'] = true;
        keyStates['w'] = true;
      } else if (vertical < 0) { // 下
        simulateKeyEvent('ArrowDown', true);
        simulateKeyEvent('s', true);
        keyStates['ArrowDown'] = true;
        keyStates['s'] = true;
      }

      if (horizontal < 0) { // 左
        simulateKeyEvent('ArrowLeft', true);
        simulateKeyEvent('a', true);
        keyStates['ArrowLeft'] = true;
        keyStates['a'] = true;
      } else if (horizontal > 0) { // 右
        simulateKeyEvent('ArrowRight', true);
        simulateKeyEvent('d', true);
        keyStates['ArrowRight'] = true;
        keyStates['d'] = true;
      }

      currentKeys.horizontal = horizontal;
      currentKeys.vertical = vertical;

      updateDirectionIndicators(horizontal, vertical);
    }

    function handleJoystickMove(touch) {
      var rect = joystickBase.getBoundingClientRect();
      var centerX = rect.left + rect.width / 2;
      var centerY = rect.top + rect.height / 2;

      var deltaX = touch.clientX - centerX;
      var deltaY = touch.clientY - centerY;

      var distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
      var maxDistance = 45;

      if (distance > maxDistance) {
        deltaX = (deltaX / distance) * maxDistance;
        deltaY = (deltaY / distance) * maxDistance;
      }

      joystickStick.style.left = (45 + deltaX) + 'px';
      joystickStick.style.top = (45 + deltaY) + 'px';

      // 计算8方向
      var angle = Math.atan2(-deltaY, deltaX) * 180 / Math.PI;
      var horizontal = 0;
      var vertical = 0;

      if (distance > 15) { // 死区
        // 8方向判断
        if (angle >= -22.5 && angle < 22.5) { // 右
          horizontal = 1;
        } else if (angle >= 22.5 && angle < 67.5) { // 右上
          horizontal = 1;
          vertical = 1;
        } else if (angle >= 67.5 && angle < 112.5) { // 上
          vertical = 1;
        } else if (angle >= 112.5 && angle < 157.5) { // 左上
          horizontal = -1;
          vertical = 1;
        } else if (angle >= 157.5 || angle < -157.5) { // 左
          horizontal = -1;
        } else if (angle >= -157.5 && angle < -112.5) { // 左下
          horizontal = -1;
          vertical = -1;
        } else if (angle >= -112.5 && angle < -67.5) { // 下
          vertical = -1;
        } else if (angle >= -67.5 && angle < -22.5) { // 右下
          horizontal = 1;
          vertical = -1;
        }
      }

      if (horizontal !== currentKeys.horizontal || vertical !== currentKeys.vertical) {
        updateKeys(horizontal, vertical);
      }
    }

    function resetJoystick() {
      joystickStick.style.left = '45px';
      joystickStick.style.top = '45px';
      updateKeys(0, 0);
    }

    joystickBase.addEventListener('touchstart', function(e) {
      e.preventDefault();
      handleJoystickMove(e.touches[0]);
    });

    joystickBase.addEventListener('touchmove', function(e) {
      e.preventDefault();
      handleJoystickMove(e.touches[0]);
    });

    joystickBase.addEventListener('touchend', function(e) {
      e.preventDefault();
      resetJoystick();
    });

    joystickBase.addEventListener('touchcancel', function(e) {
      e.preventDefault();
      resetJoystick();
    });

    console.log('Virtual joystick initialized for mobile');
  }

  // DOM 加载完成后初始化
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initJoystick);
  } else {
    initJoystick();
  }
})();
