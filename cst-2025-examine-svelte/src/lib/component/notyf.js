import { Notyf } from 'notyf';
import 'notyf/notyf.min.css';

let notyf;

const initializeNotyf = () => {
  notyf = new Notyf({
    duration: 3000,
    className: 'x-notification',
    position: { x: 'right', y: 'top' },
  });
};

if (typeof window !== 'undefined') {// 确保代码在客户端执行
  window.onload = () => {
    initializeNotyf();
  };
}

const getNotyf = () => {
  return notyf;
};

export { getNotyf };
