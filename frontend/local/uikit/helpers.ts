import { ITNScrollOptions } from "./interfaces";
import { PartialOptions } from "overlayscrollbars";

/**
 * Функция получения объекта опций для скроллбара
 *
 * @param options - настройки скроллбара
 * @returns - сформированный объект настроек
 */
export function setScrollbarsOptions(
  options?: ITNScrollOptions
): PartialOptions {
  return {
    scrollbars: {
      theme: "tn-scroll__theme",
      autoHide:
        options?.autoHide === true
          ? "scroll"
          : typeof options?.autoHide === "string"
            ? options.autoHide
            : "never",
      autoHideDelay:
        typeof options?.autoHideDelay === "number"
          ? options.autoHideDelay
          : 1300,
      autoHideSuspend: true,
      dragScroll: !options?.disabled,
      pointers: options?.disabled ? [] : ["mouse", "touch", "pen"],
      clickScroll: !options?.disabled
    }
  };
}

/**
 * Функция выполнения функции с задержкой
 *
 * @param callback - функция, вызываемая по истечении задержки
 * @param delay - длительность задержки
 */
export function debounce(callback: Function, delay: number) {
  let timeoutID: any = null;

  return function () {
    clearTimeout(timeoutID);
    timeoutID = setTimeout(() => {
      // @ts-ignore
      // eslint-disable-next-line prefer-rest-params
      callback.apply(this, arguments);
    }, delay);
  };
}

/**
 * Установка тёмной темы. Добавленный класс меняет цвета css-переменных на их цветопару тёмной темы.
 *
 * @param dark { "dark" | "light" } Тёмная тема
 */
export function setDarkTheme(dark: "dark" | "light") {
  document.body.classList.add("tn-dark-theme_no-transition");
  setTimeout(() => {
    document.body.classList.remove("tn-dark-theme_no-transition");
  }, 10);
  if (dark === "dark") {
    document.body.classList.add("tn-dark-theme");
  } else if (dark === "light") {
    document.body.classList.remove("tn-dark-theme");
  }
}
