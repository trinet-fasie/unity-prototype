export default {
  nav: {
    library: 'Библиотека',
    objects: 'Объекты',
    locations: 'Локации',
    worlds: 'Управление мирами',
    worldStructure: 'Структура мира',
    worldConfigurations: 'Конфигурации мира',
    groupLogicEditor: 'Редактор логики'
  },
  lang: {
    ru: {
      title: 'Русский',
      icon: 'flag-icon flag-icon-ru'
    },
    en: {
      title: 'English',
      icon: 'flag-icon flag-icon-us'
    }
  },
  error404: {
    message: 'Страница не найдена.',
    details: 'Странница с указанным адресом не найдена.'
  },
  objects: {
    list: 'Библиотека объектов',
    installUpdated: 'Объект "{name}" успешно обновлен',
    installCreated: 'Объект "{name}" успешно добавлен',
    deleteRestriction: 'Пока объект используется в мирах, его не возможно удалить'
  },
  locations: {
    list: 'Библиотека локаций',
    installUpdated: 'Локация "{name}" успешно обновлена',
    installCreated: 'Локация "{name}" успешно добавлена',
    deleteRestriction: 'Пока локация используется в мирах, ее не возможно удалить'
  },
  worlds: {
    list: 'Управление мирами',
    configurations: 'Конфигураций',
    deleteConfirm: 'При удалении мира будут удалены все настройки и логика. Отменить операцию будет не возможно.',
    addButton: 'Добавить мир',
    addTitle: 'Добавление мира',
    editTitle: 'Редактирование мира',
    importButton: 'Импорт мира',
    importTitle: 'Выберите файл для импорта',
    fields: {
      name: 'Имя',
      file: 'Файл'
    }
  },
  world: {
    structure: {
      header: 'Структура мира',
      locations: 'Локации',
      location: {
        deleteConfirm: 'При удалении локации буду удалены все группы, объекты и их логика. Отменить операцию будет не возможно.',
        addButton: 'Добавить локацию',
        addTitle: 'Добавление локации',
        editTitle: 'Редактирование локации',
        fields: {
          name: 'Имя',
          locationId: 'Шаблон локации'
        }
      },
      group: {
        deleteConfirm: 'При удалении группы буду удалены все объекты и их логика. Отменить операцию будет не возможно.',
        addButton: 'Добавить группу',
        addTitle: 'Добавление группы',
        editTitle: 'Редактирование группы',
        fields: {
          name: 'Имя'
        }
      },
      groupObject: {
        editTitle: 'Редактирование объекта',
        alreadyNamed: 'Объект с таким именем уже существует',
        fields: {
          name: 'Имя'
        }
      },
      configurations: 'Конфигурации',
      configuration: {
        deleteConfirm: 'При удалении конфигурации буду удалены все настройки. Отменить операцию будет не возможно.',
        addButton: 'Добавить конфигурацию',
        addTitle: 'Добавление конфигурации',
        editTitle: 'Редактирование конфигурации',
        emptyGroupsAlert: 'Для добавления конфигурации должна быть создана хотя бы одна локация и группа.',
        fields: {
          name: 'Имя',
          selectedGroups: 'Используемые группы',
          startWorldLocationId: 'Стартовая локация'
        },
        validation: {
          startWorldLocationRequired: 'Выберите используемые группы, а затем стартовую локацию'
        }
      }
    }
  },
  deleteConfirm: {
    title: 'Вы уверены?'
  },
  actions: {
    delete: 'Удалить',
    cancel: 'Отменить',
    add: 'Добавить',
    edit: 'Изменить',
    apply: 'Применить',
    applyAndOpenEditor: 'Применить и показать код',
    upload: 'Загрузить',
    close: 'Закрыть',
    reset: 'Отменить изменения',
    vr: 'Запустить в VR',
    viewInVR: 'Запустить в VR в режиме просмотра',
    editInVR: 'Запустить в VR в режиме редактирования',
    showVisualEditor: 'Открыть визуальный редактор логики',
    showCodeEditor: 'Открыть редактор кода',
    exportToFile: 'Экспорт в файл',
    buildPackage: 'Собрать пакет',
    importFromFile: 'Импорт из файла'
  },
  validation: {
    required: 'Обязательное поле',
    success: '<i class="fa fa-check"></i> Ok'
  },
  installDropzone: {
    filenamePrefix: 'Установка ',
    filesizePrefix: 'Размер: ',
    close: 'Закрыть',
    dictDefaultMessage: '<br>Переместите пакеты в зону либо выберите с диска, кликнув по зоне',
    dictCancelUpload: 'Отменить загрузку',
    dictCancelUploadConfirmation: 'Вы уверены что хотите отменить?',
    dictFallbackMessage: 'Ваш браузер не поддерживает загрузку файлов через drag and drop.',
    dictFileTooBig: 'Файл слишком большой ({{filesize}}MiB). Максимальный размер: {{maxFilesize}}MiB.',
    dictInvalidFileType: `Вы не можете загружать файлы этого типа`,
    dictMaxFilesExceeded: 'Вы не можете загружать такое количество файлов. (Максимальгное количество: {{maxFiles}})',
    dictRemoveFile: 'Удалить',
    dictRemoveFileConfirmation: null,
    dictResponseError: 'Ошибка загрузки файла. Ответ сервера: {{statusCode}}'
  },
  table: {
    filter: 'Фильтр:',
    filterClear: 'Сбросить',
    emptyText: 'Нет записей',
    emptyFilteredText: 'Нет записей удовлетворяющих фильтру',
    columns: {
      id: 'ID',
      name: 'Имя',
      type: 'Тип',
      usages: 'Использований',
      createdAt: 'Добавлено',
      updatedAt: 'Обновлено',
      tags: 'Теги',
      actions: 'Действия'
    }
  },
  editor: {
    category: {
      logic: 'Логика',
      actions: 'Действия',
      objects: 'Объекты',
      state: 'Переменные',
      lists: 'Списки',
      loops: 'Циклы',
      math: 'Математика',
      text: 'Текст',
      variables: 'Переменные',
      functions: 'Процедуры',
      events: 'События',
      other: 'Прочее'
    },
    anyObject: 'любой',
    actionBlockPrefix: 'выполнить',
    setterBlockPrefix: 'установить',
    eventsBlockPrefix: 'на событие',
    eventsBlockParams: 'параметры',
    allObjectInstancesTitle: 'список объектов типа "{name}"',
    unsavedChangesConfirm: 'Вы изменили логику и не применили изменения. Выйти с потерей изменений?'
  },
  tags: {
    addTag: 'Добавить тег'
  }
}
