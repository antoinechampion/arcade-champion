import type { Plugin } from 'vue';

function checkControllerInputs() {
    const gamepads : Gamepad[] = navigator.getGamepads();
    if (!gamepads || gamepads.length === 0) {
        requestAnimationFrame(checkControllerInputs);
        return;
    }

    const gp : Gamepad = gamepads[0];
    if (!gp || !gp.buttons) {
        requestAnimationFrame(checkControllerInputs);
        return;
    }

    if (gp.buttons.some(b => b.pressed)) {
        console.log("Button pressed");
    }



    requestAnimationFrame(checkControllerInputs);
}

export const GamepadPlugin: Plugin = {
    install(app) {
        requestAnimationFrame(checkControllerInputs);
    }
}