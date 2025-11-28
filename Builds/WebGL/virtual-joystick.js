// 虚拟摇杆控制 - 真正的圆形摇杆
(function() {
  // 检测移动设备或触摸屏
  var isMobile = /iPhone|iPad|iPod|Android/i.test(navigator.userAgent);
  var hasTouch = 'ontouchstart' in window || navigator.maxTouchPoints > 0;

  // 添加调试信息
  console.log('Virtual Joystick - Device Detection:');
  console.log('  User Agent:', navigator.userAgent);
  console.log('  Is Mobile:', isMobile);
  console.log('  Has Touch:', hasTouch);
  console.log('  Force Show (add ?joystick=1 to URL):', window.location.search.includes('joystick=1'));

  // 如果不是移动设备且没有触摸屏，并且没有强制显示参数，则退出
  if (!isMobile && !hasTouch && !window.location.search.includes('joystick=1')) {
    console.log('Virtual Joystick - Skipped (not a mobile/touch device)');
    return;
  }

  console.log('Virtual Joystick - Initializing...');

  // 等待 Unity canvas 加载
  function initJoystick() {
    var canvas = document.querySelector("#unity-canvas");
    if (!canvas) {
      console.log('Virtual Joystick - Waiting for Unity canvas...');
      setTimeout(initJoystick, 100);
      return;
    }

    console.log('Virtual Joystick - Unity canvas found, creating joystick...');

    // 创建摇杆 HTML 结构
    var joystickHTML = `
      <div id="virtual-joystick">
        <div id="joystick-base"></div>
        <div id="joystick-stick"></div>
      </div>
    `;

    document.body.insertAdjacentHTML('beforeend', joystickHTML);
    console.log('Virtual Joystick - HTML structure inserted');

    var joystick = document.getElementById('virtual-joystick');
    var joystickStick = document.getElementById('joystick-stick');
    var joystickBase = document.getElementById('joystick-base');

    if (!joystick || !joystickStick || !joystickBase) {
      console.error('Virtual joystick elements not found after insertion!');
      console.error('  joystick:', joystick);
      console.error('  joystickStick:', joystickStick);
      console.error('  joystickBase:', joystickBase);
      return;
    }

    console.log('Virtual Joystick - Elements found, setting display...');
    joystick.style.display = 'block';

    // 确保摇杆可见
    var computedStyle = window.getComputedStyle(joystick);
    console.log('Virtual Joystick - Display style:', computedStyle.display);
    console.log('Virtual Joystick - Position:', computedStyle.position);
    console.log('Virtual Joystick - Bottom:', computedStyle.bottom);
    console.log('Virtual Joystick - Left:', computedStyle.left);
    console.log('Virtual Joystick - Z-index:', computedStyle.zIndex);

    var currentInput = {
      horizontal: 0,  // -1.0 到 1.0
      vertical: 0     // -1.0 到 1.0
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

    var isDragging = false;
    var centerX = 0;
    var centerY = 0;

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
      var threshold = 0.3; // 触发阈值

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
    }

    function handleJoystickMove(clientX, clientY) {
      var rect = joystickBase.getBoundingClientRect();
      centerX = rect.left + rect.width / 2;
      centerY = rect.top + rect.height / 2;

      var deltaX = clientX - centerX;
      var deltaY = clientY - centerY;

      var distance = Math.sqrt(deltaX * deltaX + deltaY * deltaY);
      var maxDistance = 40; // 摇杆最大移动距离

      // 限制摇杆移动范围
      if (distance > maxDistance) {
        deltaX = (deltaX / distance) * maxDistance;
        deltaY = (deltaY / distance) * maxDistance;
        distance = maxDistance;
      }

      // 更新摇杆视觉位置
      joystickStick.style.left = (40 + deltaX) + 'px';
      joystickStick.style.top = (40 + deltaY) + 'px';

      // 计算归一化的输入值 (-1.0 到 1.0)
      var deadZone = 8; // 死区半径
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
      isDragging = false;
      joystickStick.style.left = '40px';
      joystickStick.style.top = '40px';
      currentInput.horizontal = 0;
      currentInput.vertical = 0;
      updateKeys();
    }

    // 触摸事件
    joystickBase.addEventListener('touchstart', function(e) {
      e.preventDefault();
      isDragging = true;
      var touch = e.touches[0];
      handleJoystickMove(touch.clientX, touch.clientY);
    });

    joystickStick.addEventListener('touchstart', function(e) {
      e.preventDefault();
      isDragging = true;
      var touch = e.touches[0];
      handleJoystickMove(touch.clientX, touch.clientY);
    });

    document.addEventListener('touchmove', function(e) {
      if (!isDragging) return;
      e.preventDefault();
      var touch = e.touches[0];
      handleJoystickMove(touch.clientX, touch.clientY);
    });

    document.addEventListener('touchend', function(e) {
      if (!isDragging) return;
      e.preventDefault();
      resetJoystick();
    });

    document.addEventListener('touchcancel', function(e) {
      if (!isDragging) return;
      e.preventDefault();
      resetJoystick();
    });

    // 鼠标事件（用于桌面测试）
    joystickBase.addEventListener('mousedown', function(e) {
      e.preventDefault();
      isDragging = true;
      handleJoystickMove(e.clientX, e.clientY);
    });

    joystickStick.addEventListener('mousedown', function(e) {
      e.preventDefault();
      isDragging = true;
      handleJoystickMove(e.clientX, e.clientY);
    });

    document.addEventListener('mousemove', function(e) {
      if (!isDragging) return;
      e.preventDefault();
      handleJoystickMove(e.clientX, e.clientY);
    });

    document.addEventListener('mouseup', function(e) {
      if (!isDragging) return;
      e.preventDefault();
      resetJoystick();
    });

    console.log('Virtual joystick initialized for mobile');
    console.log('Virtual Joystick - Setup complete! Joystick should be visible at bottom-left corner.');
  }

  // DOM 加载完成后初始化
  if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initJoystick);
  } else {
    initJoystick();
  }
})();
