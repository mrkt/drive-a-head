// 虚拟摇杆控制 - 模拟摇杆版本
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

    var currentInput = {
      horizontal: 0,  // -1.0 到 1.0 的连续值
      vertical: 0     // -1.0 到 1.0 的连续值
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

    var updateInterval = null;

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

    function updateKeys() {
      var threshold = 0.3; // 阈值，小于这个值不触发按键

      // 垂直方向
      var shouldPressUp = currentInput.vertical > threshold;
      var shouldPressDown = currentInput.vertical < -threshold;

      // 水平方向
      var shouldPressLeft = currentInput.horizontal < -threshold;
      var shouldPressRight = currentInput.horizontal > threshold;

      // 更新上/下键
      if (shouldPressUp && !keyStates['ArrowUp']) {
        simulateKeyEvent('ArrowUp', true);
        simulateKeyEvent('w', true);
        keyStates['ArrowUp'] = true;
        keyStates['w'] = true;
      } else if (!shouldPressUp && keyStates['ArrowUp']) {
        simulateKeyEvent('ArrowUp', false);
        simulateKeyEvent('w', false);
        keyStates['ArrowUp'] = false;
        keyStates['w'] = false;
      }

      if (shouldPressDown && !keyStates['ArrowDown']) {
        simulateKeyEvent('ArrowDown', true);
        simulateKeyEvent('s', true);
        keyStates['ArrowDown'] = true;
        keyStates['s'] = true;
      } else if (!shouldPressDown && keyStates['ArrowDown']) {
        simulateKeyEvent('ArrowDown', false);
        simulateKeyEvent('s', false);
        keyStates['ArrowDown'] = false;
        keyStates['s'] = false;
      }

      // 更新左/右键
      if (shouldPressLeft && !keyStates['ArrowLeft']) {
        simulateKeyEvent('ArrowLeft', true);
        simulateKeyEvent('a', true);
        keyStates['ArrowLeft'] = true;
        keyStates['a'] = true;
      } else if (!shouldPressLeft && keyStates['ArrowLeft']) {
        simulateKeyEvent('ArrowLeft', false);
        simulateKeyEvent('a', false);
        keyStates['ArrowLeft'] = false;
        keyStates['a'] = false;
      }

      if (shouldPressRight && !keyStates['ArrowRight']) {
        simulateKeyEvent('ArrowRight', true);
        simulateKeyEvent('d', true);
        keyStates['ArrowRight'] = true;
        keyStates['d'] = true;
      } else if (!shouldPressRight && keyStates['ArrowRight']) {
        simulateKeyEvent('ArrowRight', false);
        simulateKeyEvent('d', false);
        keyStates['ArrowRight'] = false;
        keyStates['d'] = false;
      }

      // 更新方向指示器
      updateDirectionIndicators();
    }

    function updateDirectionIndicators() {
      var indicators = document.querySelectorAll('.direction-indicator');
      indicators.forEach(function(ind) { ind.classList.remove('active'); });

      var threshold = 0.3;
      var h = currentInput.horizontal;
      var v = currentInput.vertical;

      if (Math.abs(h) < threshold && Math.abs(v) < threshold) {
        return; // 死区，不显示任何方向
      }

      // 计算角度来确定最接近的方向
      var angle = Math.atan2(-v, h) * 180 / Math.PI;

      if (angle >= -22.5 && angle < 22.5) {
        document.querySelector('.dir-right').classList.add('active');
      } else if (angle >= 22.5 && angle < 67.5) {
        document.querySelector('.dir-up-right').classList.add('active');
      } else if (angle >= 67.5 && angle < 112.5) {
        document.querySelector('.dir-up').classList.add('active');
      } else if (angle >= 112.5 && angle < 157.5) {
        document.querySelector('.dir-up-left').classList.add('active');
      } else if (angle >= 157.5 || angle < -157.5) {
        document.querySelector('.dir-left').classList.add('active');
      } else if (angle >= -157.5 && angle < -112.5) {
        document.querySelector('.dir-down-left').classList.add('active');
      } else if (angle >= -112.5 && angle < -67.5) {
        document.querySelector('.dir-down').classList.add('active');
      } else if (angle >= -67.5 && angle < -22.5) {
        document.querySelector('.dir-down-right').classList.add('active');
      }
    }

    function handleJoystickMove(touch) {
      var rect = joystickBase.getBoundingClientRect();
      var centerX = rect.left + rect.width / 2;
      var centerY = rect.top + rect.height / 2;

      var deltaX = touch.clientX - centerX;
      var deltaY = touch.clientY - centerY;

      var distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
      var maxDistance = 60; // 增加最大距离，让摇杆更灵敏

      // 限制摇杆移动范围
      if (distance > maxDistance) {
        deltaX = (deltaX / distance) * maxDistance;
        deltaY = (deltaY / distance) * maxDistance;
        distance = maxDistance;
      }

      // 更新摇杆视觉位置
      joystickStick.style.left = (45 + deltaX) + 'px';
      joystickStick.style.top = (45 + deltaY) + 'px';

      // 计算归一化的输入值 (-1.0 到 1.0)
      var deadZone = 10; // 死区半径
      if (distance < deadZone) {
        currentInput.horizontal = 0;
        currentInput.vertical = 0;
      } else {
        // 归一化并应用死区
        var normalizedDistance = (distance - deadZone) / (maxDistance - deadZone);
        normalizedDistance = Math.min(normalizedDistance, 1.0);

        currentInput.horizontal = (deltaX / distance) * normalizedDistance;
        currentInput.vertical = -(deltaY / distance) * normalizedDistance; // Y轴反转
      }

      updateKeys();
    }

    function resetJoystick() {
      joystickStick.style.left = '45px';
      joystickStick.style.top = '45px';
      currentInput.horizontal = 0;
      currentInput.vertical = 0;
      updateKeys();
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

    console.log('Virtual joystick (analog mode) initialized for mobile');
  }

  // DOM 加载完成后初始化
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initJoystick);
  } else {
    initJoystick();
  }
})();
