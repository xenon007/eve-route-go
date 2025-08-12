/**
 * validateSystems проверяет корректность введённых систем.
 * Возвращает ключ сообщения об ошибке или null.
 */
export function validateSystems(start: string, end: string): string | null {
  if (!start || !end) {
    return "capital.validation-required";
  }
  if (start === end) {
    return "capital.validation-same";
  }
  return null;
}
