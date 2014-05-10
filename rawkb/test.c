#include "unistd.h"
#include "linux/kd.h"
#include "termios.h"
#include "fcntl.h"
#include "sys/ioctl.h"

static struct termios tty_attr_old;
static int old_keyboard_mode;

int setupKeyboard()
{
    struct termios tty_attr;
    int flags;

    /* make stdin non-blocking */
    flags = fcntl(0, F_GETFL);
    flags |= O_NONBLOCK;
    fcntl(0, F_SETFL, flags);

    /* save old keyboard mode */
    if (ioctl(0, KDGKBMODE, &old_keyboard_mode) < 0) {
	return 0;
    }

    tcgetattr(0, &tty_attr_old);

    /* turn off buffering, echo and key processing */
    tty_attr = tty_attr_old;
    tty_attr.c_lflag &= ~(ICANON | ECHO | ISIG);
    tty_attr.c_iflag &= ~(ISTRIP | INLCR | ICRNL | IGNCR | IXON | IXOFF);
    tcsetattr(0, TCSANOW, &tty_attr);

    ioctl(0, KDSKBMODE, K_RAW);
    return 1;
}

void restoreKeyboard()
{
    tcsetattr(0, TCSAFLUSH, &tty_attr_old);
    ioctl(0, KDSKBMODE, old_keyboard_mode);
}

char readOnce()
{
    char buf[1];
    int res;

    /* read scan code from stdin */
    res = read(0, &buf[0], 1);
    return buf[0];
}

void readKeyboard()
{
    char buf[1];
    int res;

    /* read scan code from stdin */
    res = read(0, &buf[0], 1);
    /* keep reading til there's no more*/
    while (res >= 0) {
	switch (buf[0]) {
	case 0x01:
            /* escape was pressed */
            break;
        case 0x81:
            /* escape was released */
            break;
        /* process more scan code possibilities here! */
	}
	res = read(0, &buf[0], 1);
    }
}

