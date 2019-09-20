export default {
  getSwitchLocationBlock: function (Blockly, locations) {
    if (!locations.length) {
      return ''
    }

    const options = []
    locations.forEach(location => {
      options.push([location.name, location.sid])
    })

    Blockly.Blocks['switch_world_location'] = {
      init: function () {
        this.jsonInit({
          'type': 'switch_world_location',
          'message0': '%{BKY_SWITCH_WORLD_LOCATION_TITLE}',
          'args0': [
            {
              'type': 'field_dropdown',
              'name': 'location',
              'options': options
            }
          ],
          'nextStatement': null,
          'previousStatement': null,
          'colour': 100,
          'tooltip': '%{BKY_SWITCH_WORLD_LOCATION_TOOLTIP}',
          'helpUrl': ''
        })
      }
    }
    Blockly.CSharp['switch_world_location'] = block => {
      const sid = block.getFieldValue('location')

      return `Location.Load("${sid}");\n`
    }

    return `
        <block type="switch_world_location"></block>
`
  },
  getSwitchConfigurationBlock: function (Blockly, configurations) {
    if (!configurations.length) {
      return ''
    }

    const options = []
    configurations.forEach(configuration => {
      options.push([configuration.name, configuration.sid])
    })

    Blockly.Blocks['switch_world_configuration'] = {
      init: function () {
        this.jsonInit({
          'type': 'switch_world_configuration',
          'message0': '%{BKY_SWITCH_WORLD_CONFIGURATION_TITLE}',
          'args0': [
            {
              'type': 'field_dropdown',
              'name': 'configuration',
              'options': options
            }
          ],
          'nextStatement': null,
          'previousStatement': null,
          'colour': 100,
          'tooltip': '%{BKY_SWITCH_WORLD_CONFIGURATION_TOOLTIP}',
          'helpUrl': ''
        })
      }
    }
    Blockly.CSharp['switch_world_configuration'] = block => {
      const sid = block.getFieldValue('configuration')

      return `Configuration.Load("${sid}");\n`
    }

    return `
        <block type="switch_world_configuration">
        </block>
`
  },
  init: function (Blockly) {
    Blockly.Blocks['event_on_init'] = {
      init: function () {
        this.jsonInit({
          'type': 'event_on_init',
          'message0': '%{BKY_EVENT_ON_INIT_TITLE}',
          'args0': [
            {
              'type': 'input_dummy'
            },
            {
              'type': 'input_statement',
              'name': 'stmt'
            }
          ],
          'colour': '%{BKY_PROCEDURES_HUE}',
          'tooltip': '%{BKY_EVENT_ON_INIT_TOOLTIP}',
          'helpUrl': ''
        })
      }
    }
    Blockly.CSharp['event_on_init'] = (block) => {
      return Blockly.CSharp.statementToCode(block, 'stmt')
    }

    Blockly.Blocks['convert_to_string'] = {
      init: function () {
        this.jsonInit({
          'type': 'convert_to_string',
          'message0': '%{BKY_CONVERT_TO_STRING_TITLE}',
          'args0': [
            {
              'type': 'input_value',
              'name': 'value'
            }
          ],
          'colour': '%{BKY_TEXTS_HUE}',
          'output': 'String',
          'tooltip': '%{BKY_CONVERT_TO_STRING_TOOLTIP}',
          'helpUrl': ''
        })
      }
    }
    Blockly.CSharp['convert_to_string'] = (block) => {
      const value = Blockly.CSharp.valueToCode(block, 'value', Blockly.CSharp.ORDER_MEMBER) || '""'
      const code = `Utils.ConvertToString(${value})`
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }
  },
  toolbox: {
    logic: `
      <block type="controls_if"></block>
      <block type="logic_operation">
        <field name="OP">AND</field>
      </block>
      <block type="logic_compare"></block>
      <block type="logic_boolean"></block>
      <block type="logic_negate"></block>
      <block type="logic_ternary"></block>
`,
    lists: `
      <block type="get_object_hierarchic_lists">
        <value name="instance">
          <shadow type="object_any"></shadow>
        </value>  
      </block>
      <block type="lists_create_with">
        <mutation items="0"></mutation>
      </block>
      <block type="lists_create_with">
        <mutation items="3"></mutation>
      </block>
      <block type="lists_repeat">
        <value name="NUM">
          <shadow type="math_number">
            <field name="NUM">5</field>
          </shadow>
        </value>
      </block>
      <block type="lists_length"></block>
      <block type="lists_isEmpty"></block>
      <block type="lists_indexOf">
        <field name="END">FIRST</field>
        <value name="VALUE">
          <block type="variables_get">
            <field name="VAR" id="bHK7kwVgDF@qE^hd-[1T" variabletype="">list</field>
          </block>
        </value>
      </block>
      <block type="lists_getIndex">
        <mutation statement="false" at="true"></mutation>
        <field name="MODE">GET</field>
        <field name="WHERE">FROM_START</field>
        <value name="VALUE">
          <block type="variables_get">
            <field name="VAR" id="bHK7kwVgDF@qE^hd-[1T" variabletype="">list</field>
          </block>
        </value>
      </block>
      <block type="lists_setIndex">
        <mutation at="true"></mutation>
        <field name="MODE">SET</field>
        <field name="WHERE">FROM_START</field>
        <value name="LIST">
          <block type="variables_get">
            <field name="VAR" id="bHK7kwVgDF@qE^hd-[1T" variabletype="">list</field>
          </block>
        </value>
      </block>
      <block type="lists_getSublist">
        <mutation at1="true" at2="true"></mutation>
        <field name="WHERE1">FROM_START</field>
        <field name="WHERE2">FROM_START</field>
        <value name="LIST">
          <block type="variables_get">
            <field name="VAR" id="bHK7kwVgDF@qE^hd-[1T" variabletype="">list</field>
          </block>
        </value>
      </block>
      <block type="lists_split">
        <mutation mode="SPLIT"></mutation>
        <field name="MODE">SPLIT</field>
        <value name="DELIM">
          <shadow type="text">
            <field name="TEXT">,</field>
          </shadow>
        </value>
      </block>
      <block type="lists_sort">
        <field name="TYPE">NUMERIC</field>
        <field name="DIRECTION">1</field>
      </block>
`,
    loops: `
      <block type="controls_repeat_ext">
        <value name="TIMES">
          <shadow type="math_number">
            <field name="NUM">10</field>
          </shadow>
        </value>
      </block>
      <block type="controls_whileUntil">
        <field name="MODE">WHILE</field>
      </block>
      <block type="controls_for">
        <field name="VAR" id="r:16Px[7,;{FGtKFDlL?" variabletype="">i</field>
        <value name="FROM">
          <shadow type="math_number">
            <field name="NUM">1</field>
          </shadow>
        </value>
        <value name="TO">
          <shadow type="math_number">
            <field name="NUM">10</field>
          </shadow>
        </value>
        <value name="BY">
          <shadow type="math_number">
            <field name="NUM">1</field>
          </shadow>
        </value>
      </block>
      <block type="controls_forEach">
        <field name="VAR" id="j[ME31VRkAZ^SG3XPU~i" variabletype="">j</field>
      </block>
      <block type="controls_flow_statements">
        <field name="FLOW">BREAK</field>
      </block>
`,
    math: `
      <block type="math_number">
        <field name="NUM">0</field>
      </block>
      <block type="math_arithmetic">
        <field name="OP">ADD</field>
        <value name="A">
          <shadow type="math_number">
            <field name="NUM">1</field>
          </shadow>
        </value>
        <value name="B">
          <shadow type="math_number">
            <field name="NUM">1</field>
          </shadow>
        </value>
      </block>
      <block type="math_single">
        <field name="OP">ROOT</field>
        <value name="NUM">
          <shadow type="math_number">
            <field name="NUM">9</field>
          </shadow>
        </value>
      </block>
      <block type="math_trig">
        <field name="OP">SIN</field>
        <value name="NUM">
          <shadow type="math_number">
            <field name="NUM">45</field>
          </shadow>
        </value>
      </block>
      <block type="math_constant">
        <field name="CONSTANT">PI</field>
      </block>
      <block type="math_number_property">
        <mutation divisor_input="false"></mutation>
        <field name="PROPERTY">EVEN</field>
        <value name="NUMBER_TO_CHECK">
          <shadow type="math_number">
            <field name="NUM">0</field>
          </shadow>
        </value>
      </block>
      <block type="math_round">
        <field name="OP">ROUND</field>
        <value name="NUM">
          <shadow type="math_number">
            <field name="NUM">3.1</field>
          </shadow>
        </value>
      </block>
      <block type="math_random_int">
        <value name="FROM">
          <shadow type="math_number">
            <field name="NUM">1</field>
          </shadow>
        </value>
        <value name="TO">
          <shadow type="math_number">
            <field name="NUM">100</field>
          </shadow>
        </value>
      </block>
      <block type="math_random_float">
      </block>
`,
    text: `
      <block type="text">
        <field name="TEXT"></field>
      </block>
      <block type="convert_to_string">
      </block>
      <block type="text_join">
        <mutation items="2"></mutation>
      </block>
      <block type="text_append">
        <field name="VAR" id="1Sw#syDIS+J9wfcv3CYZ" variabletype="">item</field>
        <value name="TEXT">
          <shadow type="text">
            <field name="TEXT"></field>
          </shadow>
        </value>
      </block>
      <block type="text_length">
        <value name="VALUE">
          <shadow type="text">
            <field name="TEXT">abc</field>
          </shadow>
        </value>
      </block>
      <block type="text_isEmpty">
        <value name="VALUE">
          <shadow type="text">
            <field name="TEXT"></field>
          </shadow>
        </value>
      </block>
      <block type="text_indexOf">
        <field name="END">FIRST</field>
        <value name="VALUE">
          <block type="variables_get">
            <field name="VAR" id="(1VKi[HL.0V=F:ObXjzI" variabletype="">text</field>
          </block>
        </value>
        <value name="FIND">
          <shadow type="text">
            <field name="TEXT">abc</field>
          </shadow>
        </value>
      </block>
      <block type="text_charAt">
        <mutation at="true"></mutation>
        <field name="WHERE">FROM_START</field>
        <value name="VALUE">
          <block type="variables_get">
            <field name="VAR" id="(1VKi[HL.0V=F:ObXjzI" variabletype="">text</field>
          </block>
        </value>
      </block>
      <block type="text_getSubstring">
        <mutation at1="true" at2="true"></mutation>
        <field name="WHERE1">FROM_START</field>
        <field name="WHERE2">FROM_START</field>
        <value name="STRING">
          <block type="variables_get">
            <field name="VAR" id="(1VKi[HL.0V=F:ObXjzI" variabletype="">text</field>
          </block>
        </value>
      </block>
      <block type="text_changeCase">
        <field name="CASE">UPPERCASE</field>
        <value name="TEXT">
          <shadow type="text">
            <field name="TEXT">abc</field>
          </shadow>
        </value>
      </block>
      <block type="text_trim">
        <field name="MODE">BOTH</field>
        <value name="TEXT">
          <shadow type="text">
            <field name="TEXT">abc</field>
          </shadow>
        </value>
      </block>
`,
    events: `
      <block type="event_on_init">
      </block>
`
  },
  i18n: {
    ru: {
      ACTIONS_HUE: 100,

      MATH_RANDOM_INT_TITLE: 'случайное целое от %1 до %2',
      MATH_RANDOM_FLOAT_TITLE_RANDOM: 'случайное вещественное от 0 до 1',
      PROCEDURES_DEFNORETURN_TITLE: 'процедура',
      PROCEDURES_DEFRETURN_TITLE: 'процедура',
      PROCEDURES_DEFNORETURN_PROCEDURE: 'Procedure Name',
      PROCEDURES_DEFRETURN_PROCEDURE: 'Procedure Name',
      PROCEDURES_BEFORE_PARAMS: 'с параметрами: ',
      PROCEDURES_CALL_BEFORE_PARAMS: '',
      PROCEDURES_DEFNORETURN_TOOLTIP: 'Создаёт процедуру, не возвращающую значение.',
      PROCEDURES_DEFNORETURN_COMMENT: 'Описание процедуры',
      PROCEDURES_DEFRETURN_RETURN: 'вернуть',
      PROCEDURES_DEFRETURN_TOOLTIP: 'Создаёт процедуру, возвращающую значение.',
      PROCEDURES_ALLOW_STATEMENTS: 'разрешить операторы',
      PROCEDURES_DEF_DUPLICATE_WARNING: 'Предупреждение: эта процедура имеет повторяющиеся параметры.',
      PROCEDURES_CALLNORETURN_HELPURL: 'https://ru.wikipedia.org/wiki/Подпрограмма',
      PROCEDURES_CALLNORETURN_TOOLTIP: "Исполняет определённую пользователем процедуру '%1'.",
      PROCEDURES_CALLRETURN_HELPURL: 'https://ru.wikipedia.org/wiki/Подпрограмма',
      PROCEDURES_CALLRETURN_TOOLTIP: "Исполняет определённую пользователем процедуру '%1' и возвращает значение.",
      PROCEDURES_MUTATORCONTAINER_TITLE: 'параметры',
      PROCEDURES_MUTATORCONTAINER_TOOLTIP: 'Добавить, удалить или изменить порядок входных параметров для этой функции.',
      PROCEDURES_MUTATORARG_TITLE: 'имя параметра:',
      PROCEDURES_MUTATORARG_TOOLTIP: 'Добавить входной параметр в процедуру.',
      PROCEDURES_HIGHLIGHT_DEF: 'Выделить определение процедуры',
      PROCEDURES_CREATE_DO: "Создать вызов '%1'",
      PROCEDURES_IFRETURN_TOOLTIP: 'Если первое значение истинно, возвращает второе значение.',
      PROCEDURES_IFRETURN_HELPURL: 'http://c2.com/cgi/wiki?GuardClause',
      PROCEDURES_IFRETURN_WARNING: 'Предупреждение: Этот блок может использоваться только внутри определения процедуры.',

      EVENT_ON_INIT_TITLE: 'в момент инициализации %1 %2',
      EVENT_ON_INIT_TOOLTIP: 'Код в этом блоке будет выполнен один раз в момент инициализации.',

      CONVERT_TO_STRING_TITLE: 'привести к строке %1',
      CONVERT_TO_STRING_TOOLTIP: 'Приводит аргумент к строке',

      SWITCH_WORLD_LOCATION_TITLE: 'загрузить локацию %1',
      SWITCH_WORLD_LOCATION_TOOLTIP: 'Загружает указанную локацию.',

      SWITCH_WORLD_CONFIGURATION_TITLE: 'загрузить конфигурацию мира %1',
      SWITCH_WORLD_CONFIGURATION_TOOLTIP: 'Загружает указанную конфигурацию мира.',

      OBJECT_ANY_TITLE: 'объект %1',
      OBJECT_ANY_TOOLTIP: 'Объект любого типа',

      CHECK_OBJECT_TYPE_TITLE: '%1 является типом %2',
      CHECK_OBJECT_TYPE_TOOLTIP: 'Проверяет тип объекта',

      CHECK_OBJECT_STATE_TITLE: '%1 %2',
      CHECK_OBJECT_STATE_TOOLTIP: 'Проверяет состояние объекта.',
      CHECK_OBJECT_STATE_ACTIVE: 'активен',
      CHECK_OBJECT_STATE_INACTIVE: 'не активен',

      GET_OBJECT_STATE_TITLE: '%1 %2',
      GET_OBJECT_STATE_TOOLTIP: 'Получает состояние объекта.',
      GET_OBJECT_STATE_ACTIVITY: 'активность',

      GET_OBJECT_HIERARCHIC_PROPERTIES_TITLE: '%1 %2',
      GET_OBJECT_HIERARCHIC_PROPERTIES_TOOLTIP: 'Возвращает иерархические свойства объекта.',
      GET_OBJECT_HIERARCHIC_PROPERTIES_PARENT: 'родитель',

      GET_OBJECT_HIERARCHIC_LISTS_TITLE: '%1 %2',
      GET_OBJECT_HIERARCHIC_LISTS_TOOLTIP: 'Возвращает иерархические списки для объекта.',
      GET_OBJECT_HIERARCHIC_LISTS_CHILDREN: 'список дочерних объектов',
      GET_OBJECT_HIERARCHIC_LISTS_DESCENDANTS: 'список дочерних объектов на всю глубину',
      GET_OBJECT_HIERARCHIC_LISTS_ANCESTRY: 'список предков объекта',

      SET_OBJECT_STATE_TITLE: 'установить %1 %2 = %3',
      SET_OBJECT_STATE_TOOLTIP: 'Устанавливает состояние объекта.',
      SET_OBJECT_STATE_ACTIVITY: 'активность',

      OBJECT_ACTIONS_TITLE: '%1 %2',
      OBJECT_ACTIONS_TOOLTIP: 'Запускает действие объекта.',
      OBJECT_ACTIONS_OPTION_ACTIVATE: 'активировать',
      OBJECT_ACTIONS_OPTION_DEACTIVATE: 'декативировать'
    },
    en: {
      ACTIONS_HUE: 100,

      EVENT_ON_INIT_TITLE: 'on init %1 %2',
      EVENT_ON_INIT_TOOLTIP: 'The code in this block will be executed once at the time of initialization.',

      CONVERT_TO_STRING_TITLE: 'to string %1',
      CONVERT_TO_STRING_TOOLTIP: 'Convert argument to string',

      SWITCH_WORLD_LOCATION_TITLE: 'load location %1',
      SWITCH_WORLD_LOCATION_TOOLTIP: 'Loads the specified location.',

      SWITCH_WORLD_CONFIGURATION_TITLE: 'load world configuration %1',
      SWITCH_WORLD_CONFIGURATION_TOOLTIP: 'Loads the specified world configuration.',

      OBJECT_ANY_TITLE: 'object %1',
      OBJECT_ANY_TOOLTIP: 'Object of any type',

      CHECK_OBJECT_TYPE_TITLE: '%1 is %2',
      CHECK_OBJECT_TYPE_TOOLTIP: 'Check object type',

      CHECK_OBJECT_STATE_TITLE: '%1 %2',
      CHECK_OBJECT_STATE_TOOLTIP: 'Check object state.',
      CHECK_OBJECT_STATE_ACTIVE: 'active',
      CHECK_OBJECT_STATE_INACTIVE: 'inactive',

      GET_OBJECT_STATE_TITLE: '%1 %2',
      GET_OBJECT_STATE_TOOLTIP: 'Get object state.',
      GET_OBJECT_STATE_ACTIVITY: 'activity',

      GET_OBJECT_HIERARCHIC_PROPERTIES_TITLE: '%1 %2',
      GET_OBJECT_HIERARCHIC_PROPERTIES_TOOLTIP: 'Returns hierarchic properties for object.',
      GET_OBJECT_HIERARCHIC_PROPERTIES_PARENT: 'parent',

      GET_OBJECT_HIERARCHIC_LISTS_TITLE: '%1 %2',
      GET_OBJECT_HIERARCHIC_LISTS_TOOLTIP: 'Returns hierarchic list for object.',
      GET_OBJECT_HIERARCHIC_LISTS_CHILDREN: 'list of children',
      GET_OBJECT_HIERARCHIC_LISTS_DESCENDANTS: 'list of descendants',
      GET_OBJECT_HIERARCHIC_LISTS_ANCESTRY: 'list of ancestry',

      SET_OBJECT_STATE_TITLE: 'set %1 %2 = %3',
      SET_OBJECT_STATE_TOOLTIP: 'Set object state.',
      SET_OBJECT_STATE_ACTIVITY: 'activity',

      OBJECT_ACTIONS_TITLE: '%1 %2',
      OBJECT_ACTIONS_TOOLTIP: 'Execute some object action.',
      OBJECT_ACTIONS_OPTION_ACTIVATE: 'activate',
      OBJECT_ACTIONS_OPTION_DEACTIVATE: 'deactivate'
    }
  }
}
