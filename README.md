# Модификатор текста
Эта утилита создана для копирования файлов (аналог `unix` команды `dd`),
с помощью которой можно прочитать данные из одного файла (или `stdin`),
совершить необходимые преобразования и записать их в новый файл (или распечатать в `stdout`).

Проект реализован в рамках курса по Go от Т-Банка

## Поддерживаемые флаги
* `-from` - путь к исходному файлу. по умолчанию, если `-from` не задан, в качестве input'а берется `stdin`;
* `-to` - путь к копии, по умолчанию используется `stdout`;
* `-offset` - количество байт внутри input'а, которое необходимо пропустить при копировании;
* `-limit` - максимальное количество читаемых байт. По умолчанию копируем все содержимое начиная с `-offset`;
* `-block-size` - размер одного блока в байтах при чтении и записи. То есть за один раз нельзя ни читать, ни записывать больше байт, чем `-block-size`;
* `-conv` - одно или несколько из возможных преобразований над текстом, разделенные запятой. Преобразования применяются после `-offset` и `-limit`.
  возможные значения параметра:
    - `upper_case` - преобразование всего текста к верхнему регистру;
    - `lower_case` - преобразование всего текста к нижнему регистру.
      если указаны и `lower_case`, и `upper_case` - ошибка;
    - `trim_spaces` - обрезание пробельных символов с начала и конца текста.

## Особенности реализации

* допустимо использовать `-limit` больше, чем размер файла. Копируется весь исходный файл до его `EOF`;
* идет проверка, что `-offset` меньше, чем размер файла. Если это не так - возвращаем ошибку;
* все ошибки пишутся в `stderr`;
* предполагается, что в качестве input'ов будут поступать данные в формате `UTF8`.



