export default {
  init: function (Blockly) {
    Blockly.CSharp = new Blockly.Generator('CSharp')

    Blockly.CSharp.addReservedWords(
      // http://msdn.microsoft.com/en-us/library/x53a06bb.aspx
      'abstract,as,base,bool,break,byte,case,catch,char,checked,class,const,continue,decimal,default,delegate,do,double,else,enum,event,explicit,extern,false,finally,fixed,float,for,foreach,goto,if,implicit,in,int,interface,internal,is,lock,long,namespace,new,null,object,operator,out,override,params,private,protected,public,readonly,ref,return,sbyte,sealed,short,sizeof,stackalloc,static,string,struct,switch,this,throw,true,try,typeof,uint,ulong,unchecked,unsafe,ushort,using,virtual,void,volatile,while'
    )

    Blockly.CSharp.ORDER_ATOMIC = 0 // 0 ""
    Blockly.CSharp.ORDER_MEMBER = 1 // . []
    Blockly.CSharp.ORDER_NEW = 1 // new
    Blockly.CSharp.ORDER_TYPEOF = 1 // typeof
    Blockly.CSharp.ORDER_FUNCTION_CALL = 1 // ()
    Blockly.CSharp.ORDER_INCREMENT = 1 // ++
    Blockly.CSharp.ORDER_DECREMENT = 1 // --
    Blockly.CSharp.ORDER_LOGICAL_NOT = 2 // !
    Blockly.CSharp.ORDER_BITWISE_NOT = 2 // ~
    Blockly.CSharp.ORDER_UNARY_PLUS = 2 // +
    Blockly.CSharp.ORDER_UNARY_NEGATION = 2 // -
    Blockly.CSharp.ORDER_MULTIPLICATION = 3 // *
    Blockly.CSharp.ORDER_DIVISION = 3 // /
    Blockly.CSharp.ORDER_MODULUS = 3 // %
    Blockly.CSharp.ORDER_ADDITION = 4 // +
    Blockly.CSharp.ORDER_SUBTRACTION = 4 // -
    Blockly.CSharp.ORDER_BITWISE_SHIFT = 5 // << >>
    Blockly.CSharp.ORDER_RELATIONAL = 6 // < <= > >=
    Blockly.CSharp.ORDER_EQUALITY = 7 // == !=
    Blockly.CSharp.ORDER_BITWISE_AND = 8 // &
    Blockly.CSharp.ORDER_BITWISE_XOR = 9 // ^
    Blockly.CSharp.ORDER_BITWISE_OR = 10 // |
    Blockly.CSharp.ORDER_LOGICAL_AND = 11 // &&
    Blockly.CSharp.ORDER_LOGICAL_OR = 12 // ||
    Blockly.CSharp.ORDER_CONDITIONAL = 13 // ?:
    Blockly.CSharp.ORDER_ASSIGNMENT = 14 // = += -= *= /= %= <<= >>= ...
    Blockly.CSharp.ORDER_COMMA = 15 // ,
    Blockly.CSharp.ORDER_NONE = 99 // (...)

    /**
     * Arbitrary code to inject into locations that risk causing infinite loops.
     * Any instances of '%1' will be replaced by the block ID that failed.
     * E.g. '  checkTimeout(%1)\n'
     * @type ?string
     */
    Blockly.CSharp.INFINITE_LOOP_TRAP = null

    Blockly.CSharp.init = function (workspace) {
      Blockly.CSharp.definitions_ = {}

      if (Blockly.Variables) {
        if (!Blockly.CSharp.variableDB_) {
          Blockly.CSharp.variableDB_ = new Blockly.Names(Blockly.CSharp.RESERVED_WORDS_)
        } else {
          Blockly.CSharp.variableDB_.reset()
        }

        Blockly.CSharp.variableDB_.setVariableMap(workspace.getVariableMap())

        const defvars = []
        // Add developer variables (not created or named by the user).
        const devVarList = Blockly.Variables.allDeveloperVariables(workspace)
        for (let i = 0; i < devVarList.length; i++) {
          defvars.push(Blockly.CSharp.variableDB_.getName(devVarList[i], Blockly.Names.DEVELOPER_VARIABLE_TYPE))
        }

        // Add user variables, but only ones that are being used.
        const variables = Blockly.Variables.allUsedVarModels(workspace)
        for (let i = 0; i < variables.length; i++) {
          defvars.push(Blockly.CSharp.variableDB_.getName(variables[i].getId(), Blockly.VARIABLE_CATEGORY_NAME))
        }

        // Declare all of the variables.
        if (defvars.length) {
          Blockly.CSharp.definitions_['variables'] = 'private dynamic ' + defvars.join(', ') + ';'
        }

        Blockly.CSharp.using = {}
        Blockly.CSharp.lockedInstancesForDelete = {}
      }
    }

    Blockly.CSharp.use = function (namespace) {
      Blockly.CSharp.using[namespace] = 1
    }

    Blockly.CSharp.lockInstance = function (instanceId) {
      Blockly.CSharp.lockedInstancesForDelete[instanceId] = 1
    }

    /* Prepend the generated code with the variable definitions. */
    Blockly.CSharp.finish = function (code) {
      let definitions = []
      for (let name in Blockly.CSharp.definitions_) {
        if (Blockly.CSharp.definitions_.hasOwnProperty(name)) {
          definitions.push(Blockly.CSharp.definitions_[name])
        }
      }
      return definitions.join('\n\n') + '\n\n\n' + code
    }

    /**
     * Naked values are top-level blocks with outputs that aren't plugged into
     * anything.  A trailing semicolon is needed to make this legal.
     * @param {string} line Line of generated code.
     * @return {string} Legal line of code.
     */
    Blockly.CSharp.scrubNakedValue = function (line) {
      return line + ';\n'
    }

    Blockly.CSharp.quote_ = function (val) {
      return Blockly.goog.string.quote(val)
    }

    /**
     * Common tasks for generating cSharp from blocks.
     * Handles comments for the specified block and any connected value blocks.
     * Calls any statements following this block.
     * @param {!Blockly.Block} block The current block.
     * @param {string} code The cSharp code created for this block.
     * @return {string} cSharp code with comments and subsequent blocks added.
     * @private
     */
    Blockly.CSharp.scrub_ = function (block, code) {
      if (code === null) {
        // Block has handled code generation itself.
        return ''
      }
      let commentCode = ''
      // Only collect comments for blocks that aren't inline.
      if (!block.outputConnection || !block.outputConnection.targetConnection) {
        // Collect comment for this block.
        let comment = block.getCommentText()
        if (comment) {
          commentCode += this.prefixLines(comment, '// ') + '\n'
        }
        // Collect comments for all value arguments.
        // Don't collect comments for nested statements.
        for (let x = 0; x < block.inputList.length; x++) {
          if (block.inputList[x].type === Blockly.INPUT_VALUE) {
            let childBlock = block.inputList[x].connection.targetBlock()
            if (childBlock) {
              let comment = this.allNestedComments(childBlock)
              if (comment) {
                commentCode += this.prefixLines(comment, '// ')
              }
            }
          }
        }
      }
      let nextBlock = block.nextConnection && block.nextConnection.targetBlock()
      let nextCode = this.blockToCode(nextBlock)
      return commentCode + code + nextCode
    }

    this.initColour(Blockly)
    this.initLists(Blockly)
    this.initLogic(Blockly)
    this.initLoops(Blockly)
    this.initMath(Blockly)
    this.initProcedures(Blockly)
    this.initText(Blockly)
    this.initVariables(Blockly)

    return Blockly.CSharp
  },

  initColour: function (Blockly) {
    Blockly.CSharp.colour = {}

    Blockly.CSharp.colour_picker = function () {
      // Colour picker.
      const code = 'ColorTranslator.FromHtml("' + this.getFieldValue('COLOUR') + '")'
      return [code, Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp.colour_random = function () {
      // Generate a random colour.
      if (!Blockly.CSharp.definitions_['colour_random']) {
        const functionName = Blockly.CSharp.variableDB_.getDistinctName('colour_random', Blockly.Generator.VARIABLE_CATEGORY_NAME)
        Blockly.CSharp.colour_random.functionName = functionName
        let func = []
        func.push('var ' + functionName + ' = new Func<Color>(() => {')
        func.push('  var random = new Random();')
        func.push('  var res = Color.FromArgb(1, random.Next(256), random.Next(256), random.Next(256));')
        func.push('  return res;')
        func.push('});')
        Blockly.CSharp.definitions_['colour_random'] = func.join('\n')
      }
      const code = Blockly.CSharp.colour_random.functionName + '()'
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.colour_rgb = function () {
      // Compose a colour from RGB components expressed as percentages.
      const red = Blockly.CSharp.valueToCode(this, 'RED', Blockly.CSharp.ORDER_COMMA) || 0
      const green = Blockly.CSharp.valueToCode(this, 'GREEN', Blockly.CSharp.ORDER_COMMA) || 0
      const blue = Blockly.CSharp.valueToCode(this, 'BLUE', Blockly.CSharp.ORDER_COMMA) || 0

      if (!Blockly.CSharp.definitions_['colour_rgb']) {
        const functionName = Blockly.CSharp.variableDB_.getDistinctName('colour_rgb', Blockly.Generator.VARIABLE_CATEGORY_NAME)
        Blockly.CSharp.colour_rgb.functionName = functionName
        const func = []
        func.push('var ' + functionName + ' = new Func<dynamic, dynamic, dynamic, Color>((r, g, b) => {')
        func.push('  r = (int)Math.Round(Math.Max(Math.Min((int)r, 100), 0) * 2.55);')
        func.push('  g = (int)Math.Round(Math.Max(Math.Min((int)g, 100), 0) * 2.55);')
        func.push('  b = (int)Math.Round(Math.Max(Math.Min((int)b, 100), 0) * 2.55);')
        func.push('  var res = Color.FromArgb(1, r, g, b);')
        func.push('  return res;')
        func.push('});')
        Blockly.CSharp.definitions_['colour_rgb'] = func.join('\n')
      }
      const code = Blockly.CSharp.colour_rgb.functionName + '(' + red + ', ' + green + ', ' + blue + ')'
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.colour_blend = function () {
      // Blend two colours together.
      const c1 = Blockly.CSharp.valueToCode(this, 'COLOUR1', Blockly.CSharp.ORDER_COMMA) || 'Color.Black'
      const c2 = Blockly.CSharp.valueToCode(this, 'COLOUR2', Blockly.CSharp.ORDER_COMMA) || 'Color.Black'
      const ratio = Blockly.CSharp.valueToCode(this, 'RATIO', Blockly.CSharp.ORDER_COMMA) || 0.5

      if (!Blockly.CSharp.definitions_['colour_blend']) {
        const functionName = Blockly.CSharp.variableDB_.getDistinctName('colour_blend', Blockly.Generator.VARIABLE_CATEGORY_NAME)
        Blockly.CSharp.colour_blend.functionName = functionName
        const func = []
        func.push('var ' + functionName + ' = new Func<Color, Color, double, Color>((c1, c2, ratio) => {')
        func.push('  ratio = Math.Max(Math.Min((double)ratio, 1), 0);')
        func.push('  var r = (int)Math.Round(c1.R * (1 - ratio) + c2.R * ratio);')
        func.push('  var g = (int)Math.Round(c1.G * (1 - ratio) + c2.G * ratio);')
        func.push('  var b = (int)Math.Round(c1.B * (1 - ratio) + c2.B * ratio);')
        func.push('  var res = Color.FromArgb(1, r, g, b);')
        func.push('  return res;')
        func.push('});')
        Blockly.CSharp.definitions_['colour_blend'] = func.join('\n')
      }
      const code = Blockly.CSharp.colour_blend.functionName + '(' + c1 + ', ' + c2 + ', ' + ratio + ')'
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }
  },
  initLists: function (Blockly) {
    Blockly.CSharp.lists = {}

    Blockly.CSharp.lists_create_empty = function () {
      return ['null', Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp.lists_create_with = function () {
      // Create a list with any number of elements of any type.
      let code = new Array(this.itemCount_)
      for (let n = 0; n < this.itemCount_; n++) {
        code[n] = Blockly.CSharp.valueToCode(this, 'ADD' + n, Blockly.CSharp.ORDER_COMMA) || 'null'
      }
      code = 'new List<dynamic> {' + code.join(', ') + '}'
      return [code, Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp.lists_repeat = function () {
      // Create a list with one element repeated.
      if (!Blockly.CSharp.definitions_['lists_repeat']) {
        // Function copied from Closure's goog.array.repeat.
        const functionName = Blockly.CSharp.variableDB_.getDistinctName('lists_repeat', Blockly.Generator.VARIABLE_CATEGORY_NAME)
        Blockly.CSharp.lists_repeat.repeat = functionName
        const func = []
        func.push('var ' + functionName + ' = new Func<dynamic, dynamic, List<dynamic>>((value, n) => {')
        func.push('  var array = new List<dynamic>(n);')
        func.push('  for (var i = 0; i < n; i++) {')
        func.push('    array.Add(value);')
        func.push('  }')
        func.push('  return array;')
        func.push('});')
        Blockly.CSharp.definitions_['lists_repeat'] = func.join('\n')
      }
      const argument0 = Blockly.CSharp.valueToCode(this, 'ITEM', Blockly.CSharp.ORDER_COMMA) || 'null'
      const argument1 = Blockly.CSharp.valueToCode(this, 'NUM', Blockly.CSharp.ORDER_COMMA) || '0'
      const code = Blockly.CSharp.lists_repeat.repeat + '(' + argument0 + ', ' + argument1 + ')'
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.lists_length = function () {
      // List length.
      let argument0 = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_FUNCTION_CALL) || 'null'
      return [argument0 + '.Count', Blockly.CSharp.ORDER_MEMBER]
    }

    Blockly.CSharp.lists_isEmpty = function () {
      // Is the list empty?
      const argument0 = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_MEMBER) || 'null'
      return [argument0 + '.Count == 0', Blockly.CSharp.ORDER_LOGICAL_NOT]
    }

    Blockly.CSharp.lists_indexOf = function () {
      // Find an item in the list.
      const operator = this.getFieldValue('END') === 'FIRST' ? 'IndexOf' : 'LastIndexOf'
      const argument0 = Blockly.CSharp.valueToCode(this, 'FIND', Blockly.CSharp.ORDER_NONE) || 'null'
      const argument1 = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_MEMBER) || 'null'
      const code = argument1 + '.' + operator + '(' + argument0 + ') + 1'
      return [code, Blockly.CSharp.ORDER_MEMBER]
    }

    Blockly.CSharp.lists_getIndex = function () {
      const mode = this.getFieldValue('MODE') || 'GET'
      const where = this.getFieldValue('WHERE') || 'FROM_START'
      let at = Blockly.CSharp.valueToCode(this, 'AT', Blockly.CSharp.ORDER_UNARY_NEGATION) || '1'
      let list = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_MEMBER) || 'null'

      if (mode === 'GET_REMOVE') {
        if (where === 'FIRST') {
          at = 1
        } else if (where === 'LAST') {
          at = list + '.Count - 1'
        } else {
          // Blockly uses one-based indicies.
          if (Blockly.isNumber(at)) {
            // If the index is a naked number, decrement it right now.
            at = parseFloat(at) - 1
          } else {
            // If the index is expression, decrement it in code.
            at = '(' + at + ' - 1)'
          }
        }

        if (where === 'FROM_END') {
          at = '(' + list + '.Count) - ' + (at + 1)
        }

        if (!Blockly.CSharp.definitions_['lists_get_remove_at']) {
          const functionName = Blockly.CSharp.variableDB_.getDistinctName('lists_get_remove_at', Blockly.Generator.VARIABLE_CATEGORY_NAME)
          Blockly.CSharp.lists_getIndex.lists_get_remove_at = functionName
          const func = []
          func.push('var ' + functionName + ' = new Func<List<dynamic>, int, dynamic>((list, index) => {')
          func.push('  var res = list[index];')
          func.push('  list.RemoveAt(index);')
          func.push('  return res;')
          func.push('});')
          Blockly.CSharp.definitions_['lists_get_remove_at'] = func.join('\n')
        }
        const code = Blockly.CSharp.lists_getIndex.lists_get_remove_at + '(' + list + ', ' + at + ')'
        return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      }

      if (where === 'FIRST') {
        if (mode === 'GET') {
          const code = list + '.First()'
          return [code, Blockly.CSharp.ORDER_MEMBER]
        } else if (mode === 'REMOVE') {
          return list + '.RemoveAt(0);\n'
        }
      } else if (where === 'LAST') {
        if (mode === 'GET') {
          const code = list + '.Last()'
          return [code, Blockly.CSharp.ORDER_MEMBER]
        } else if (mode === 'REMOVE') {
          return list + '.RemoveAt(' + list + '.Count - 1);\n'
        }
      } else if (where === 'FROM_START') {
        if (mode === 'GET') {
          const code = list + '[' + at + ']'
          return [code, Blockly.CSharp.ORDER_MEMBER]
        } else if (mode === 'REMOVE') {
          return list + '.RemoveAt(' + at + ');\n'
        }
      } else if (where === 'FROM_END') {
        if (mode === 'GET') {
          const code = list + '[list.Count - ' + at + ']'
          return [code, Blockly.CSharp.ORDER_MEMBER]
        } else if (mode === 'REMOVE') {
          return list + '.RemoveAt(list.Count - ' + at + ');\n'
        }
      } else if (where === 'RANDOM') {
        if (!Blockly.CSharp.definitions_['lists_get_random_item']) {
          const functionName = Blockly.CSharp.variableDB_.getDistinctName('lists_get_random_item', Blockly.Generator.VARIABLE_CATEGORY_NAME)
          Blockly.CSharp.lists_getIndex.random = functionName
          const func = []
          func.push('var ' + functionName + ' = new Func<List<dynamic>, bool, dynamic>((list, remove) => {')
          func.push('  var x = (new Random()).Next(list.Count);')
          func.push('  if (remove) {')
          func.push('    var res = list[x];')
          func.push('    list.RemoveAt(x);')
          func.push('    return res;')
          func.push('  } else {')
          func.push('    return list[x];')
          func.push('  }')
          func.push('});')
          Blockly.CSharp.definitions_['lists_get_random_item'] = func.join('\n')
        }
        const code = Blockly.CSharp.lists_getIndex.random + '(' + list + ', ' + (mode !== 'GET') + ')'
        if (mode === 'GET') {
          return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
        } else if (mode === 'REMOVE') {
          return code + ';\n'
        }
      }
      throw Error('Unhandled combination (lists_getIndex).')
    }

    Blockly.CSharp.lists_setIndex = function () {
      // Set element at index.
      const list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_MEMBER) || 'null'
      const mode = this.getFieldValue('MODE') || 'GET'
      const where = this.getFieldValue('WHERE') || 'FROM_START'
      let at = Blockly.CSharp.valueToCode(this, 'AT', Blockly.CSharp.ORDER_NONE) || '1'
      const value = Blockly.CSharp.valueToCode(this, 'TO', Blockly.CSharp.ORDER_ASSIGNMENT) || 'null'

      if (where === 'FIRST') {
        if (mode === 'SET') {
          return list + '[0] = ' + value + ';\n'
        } else if (mode === 'INSERT') {
          return list + '.Insert(0, ' + value + ');\n'
        }
      } else if (where === 'LAST') {
        if (mode === 'SET') {
          return list + '[' + list + '.Count - 1] = ' + value + ';\n'
        } else if (mode === 'INSERT') {
          return list + '.Add(' + value + ');\n'
        }
      } else if (where === 'FROM_START') {
        // Blockly uses one-based indicies.
        if (Blockly.isNumber(at)) {
          // If the index is a naked number, decrement it right now.
          at = parseFloat(at) - 1
        } else {
          // If the index is dynamic, decrement it in code.
          at = '(' + list + '.Count) - ' + (at + 1)
        }
        if (mode === 'SET') {
          return list + '[' + at + '] = ' + value + ';\n'
        } else if (mode === 'INSERT') {
          return list + '.Insert(' + at + ', ' + value + ');\n'
        }
      } else if (where === 'FROM_END') {
        if (mode === 'SET') {
          return list + '[' + list + '.Count - ' + at + '] = ' + value + ';\n'
        } else if (mode === 'INSERT') {
          return list + '.Insert(' + list + '.Count - ' + at + ', ' + value + ');\n'
        }
      } else if (where === 'RANDOM') {
        const xVar = Blockly.CSharp.variableDB_.getDistinctName('tmp_x', Blockly.VARIABLE_CATEGORY_NAME)
        const code = 'var ' + xVar + ' = (new Random()).Next(' + list + '.Count);\n'
        if (mode === 'SET') {
          return code + list + '[' + xVar + '] = ' + value + ';\n'
        } else if (mode === 'INSERT') {
          return code + list + '.Insert(' + xVar + ', ' + value + ');\n'
        }
      }
      throw Error('Unhandled combination (lists_setIndex).')
    }

    Blockly.CSharp.lists_getSublist = function () {
      // Get sublist.
      const list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_MEMBER) || 'null'
      const where1 = this.getFieldValue('WHERE1')
      const where2 = this.getFieldValue('WHERE2')
      let at1 = Blockly.CSharp.valueToCode(this, 'AT1', Blockly.CSharp.ORDER_NONE) || '1'
      let at2 = Blockly.CSharp.valueToCode(this, 'AT2', Blockly.CSharp.ORDER_NONE) || '1'
      let code = ''
      if (where1 === 'FIRST' && where2 === 'LAST') {
        code = 'new List<dynamic>(' + list + ')'
      } else {
        if (!Blockly.CSharp.definitions_['lists_get_sublist']) {
          const functionName = Blockly.CSharp.variableDB_.getDistinctName('lists_get_sublist', Blockly.Generator.VARIABLE_CATEGORY_NAME)
          Blockly.CSharp.lists_getSublist.func = functionName
          const func = []
          func.push('var ' + functionName + ' = new Func<List<dynamic>, dynamic, int, dynamic, int, List<dynamic>>((list, where1, at1, where2, at2) => {')
          func.push('  var getIndex = new Func<dynamic, int, int>((where, at) => {')
          func.push('    if (where == "FROM_START") {')
          func.push('      at--;')
          func.push('    } else if (where == "FROM_END") {')
          func.push('      at = list.Count - at;')
          func.push('    } else if (where == "FIRST") {')
          func.push('      at = 0;')
          func.push('    } else if (where == "LAST") {')
          func.push('      at = list.Count - 1;')
          func.push('    } else {')
          func.push('      throw new ApplicationException("Unhandled option (lists_getSublist).");')
          func.push('    }')
          func.push('    return at;')
          func.push('  });')
          func.push('  at1 = getIndex(where1, at1);')
          func.push('  at2 = getIndex(where2, at2);')
          func.push('  return list.GetRange(at1, at2 - at1 + 1);')
          func.push('});')
          Blockly.CSharp.definitions_['lists_get_sublist'] = func.join('\n')
        }
        code = Blockly.CSharp.lists_getSublist.func + '(' + list + ', "' + where1 + '", ' + at1 + ', "' + where2 + '", ' + at2 + ')'
      }
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }
  },
  initLogic: function (Blockly) {
    Blockly.CSharp.logic = {}

    Blockly.CSharp.controls_if = function () {
      // If/elseif/else condition.
      let n = 0
      let argument = Blockly.CSharp.valueToCode(this, 'IF' + n, Blockly.CSharp.ORDER_NONE) || 'false'
      let branch = Blockly.CSharp.statementToCode(this, 'DO' + n)
      let code = 'if (' + argument + ') {\n' + branch + '}'
      for (n = 1; n <= this.elseifCount_; n++) {
        argument = Blockly.CSharp.valueToCode(this, 'IF' + n, Blockly.CSharp.ORDER_NONE) || 'false'
        branch = Blockly.CSharp.statementToCode(this, 'DO' + n)
        code += ' else if (' + argument + ') {\n' + branch + '}\n'
      }
      if (this.elseCount_) {
        branch = Blockly.CSharp.statementToCode(this, 'ELSE')
        code += ' else {\n' + branch + '}\n'
      }
      return code + '\n'
    }

    Blockly.CSharp.logic_compare = function () {
      // Comparison operator.
      const mode = this.getFieldValue('OP')
      const operator = Blockly.CSharp.logic_compare.OPERATORS[mode]
      const order = (operator === '==' || operator === '!=') ? Blockly.CSharp.ORDER_EQUALITY : Blockly.CSharp.ORDER_RELATIONAL
      const argument0 = Blockly.CSharp.valueToCode(this, 'A', order) || 'null'
      const argument1 = Blockly.CSharp.valueToCode(this, 'B', order) || 'null'
      const code = argument0 + ' ' + operator + ' ' + argument1
      return [code, order]
    }

    Blockly.CSharp.logic_compare.OPERATORS = {
      EQ: '==',
      NEQ: '!=',
      LT: '<',
      LTE: '<=',
      GT: '>',
      GTE: '>='
    }

    Blockly.CSharp.logic_operation = function () {
      // Operations 'and', 'or'.
      const operator = (this.getFieldValue('OP') === 'AND') ? '&&' : '||'
      const order = (operator === '&&') ? Blockly.CSharp.ORDER_LOGICAL_AND : Blockly.CSharp.ORDER_LOGICAL_OR
      const argument0 = Blockly.CSharp.valueToCode(this, 'A', order) || 'false'
      const argument1 = Blockly.CSharp.valueToCode(this, 'B', order) || 'false'
      const code = argument0 + ' ' + operator + ' ' + argument1
      return [code, order]
    }

    Blockly.CSharp.logic_negate = function () {
      // Negation.
      const order = Blockly.CSharp.ORDER_LOGICAL_NOT
      const argument0 = Blockly.CSharp.valueToCode(this, 'BOOL', order) || 'false'
      const code = '!' + argument0
      return [code, order]
    }

    Blockly.CSharp.logic_boolean = function () {
      // Boolean values true and false.
      const code = (this.getFieldValue('BOOL') === 'TRUE') ? 'true' : 'false'
      return [code, Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp.logic_null = function () {
      // Null data type.
      return ['null', Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp.logic_ternary = function () {
      // Ternary operator.
      const valueIf = Blockly.CSharp.valueToCode(this, 'IF', Blockly.CSharp.ORDER_CONDITIONAL) || 'false'
      const valueThen = Blockly.CSharp.valueToCode(this, 'THEN', Blockly.CSharp.ORDER_CONDITIONAL) || 'null'
      const valueElse = Blockly.CSharp.valueToCode(this, 'ELSE', Blockly.CSharp.ORDER_CONDITIONAL) || 'null'
      const code = valueIf + ' ? ' + valueThen + ' : ' + valueElse
      return [code, Blockly.CSharp.ORDER_CONDITIONAL]
    }
  },
  initLoops: function (Blockly) {
    Blockly.CSharp.control = {}

    Blockly.CSharp.controls_repeat = function () {
      // Repeat n times (internal number).
      const repeats = Number(this.getFieldValue('TIMES'))
      let branch = Blockly.CSharp.statementToCode(this, 'DO')
      if (Blockly.CSharp.INFINITE_LOOP_TRAP) {
        branch = Blockly.CSharp.INFINITE_LOOP_TRAP.replace(/%1/g, '\'' + this.id + '\'') + branch
      }
      const loopVar = Blockly.CSharp.variableDB_.getDistinctName('count', Blockly.VARIABLE_CATEGORY_NAME)
      return 'for (var ' + loopVar + ' = 0; ' +
        loopVar + ' < ' + repeats + '; ' +
        loopVar + '++) {\n' +
        branch + '}\n'
    }

    Blockly.CSharp.controls_repeat_ext = function () {
      // Repeat n times (external number).
      const repeats = Blockly.CSharp.valueToCode(this, 'TIMES', Blockly.CSharp.ORDER_ASSIGNMENT) || '0'
      let branch = Blockly.CSharp.statementToCode(this, 'DO')
      if (Blockly.CSharp.INFINITE_LOOP_TRAP) {
        branch = Blockly.CSharp.INFINITE_LOOP_TRAP.replace(/%1/g, '\'' + this.id + '\'') + branch
      }
      let code = ''
      const loopVar = Blockly.CSharp.variableDB_.getDistinctName('count', Blockly.VARIABLE_CATEGORY_NAME)
      let endVar = repeats
      if (!repeats.match(/^\w+$/) && !Blockly.isNumber(repeats)) {
        endVar = Blockly.CSharp.variableDB_.getDistinctName('repeat_end', Blockly.VARIABLE_CATEGORY_NAME)
        code += 'var ' + endVar + ' = ' + repeats + ';\n'
      }
      code += 'for (var ' + loopVar + ' = 0; ' + loopVar + ' < ' + endVar + '; ' + loopVar + '++) {\n' + branch + '}\n'
      return code
    }

    Blockly.CSharp.controls_whileUntil = function () {
      // Do while/until loop.
      const until = this.getFieldValue('MODE') === 'UNTIL'
      let argument0 = Blockly.CSharp.valueToCode(this, 'BOOL', until ? Blockly.CSharp.ORDER_LOGICAL_NOT : Blockly.CSharp.ORDER_NONE) || 'false'
      let branch = Blockly.CSharp.statementToCode(this, 'DO')
      if (Blockly.CSharp.INFINITE_LOOP_TRAP) {
        branch = Blockly.CSharp.INFINITE_LOOP_TRAP.replace(/%1/g, '\'' + this.id + '\'') + branch
      }
      if (until) {
        argument0 = '!' + argument0
      }
      return 'while (' + argument0 + ') {\n' + branch + '}\n'
    }

    Blockly.CSharp.controls_for = function () {
      // For loop.
      const variable0 = Blockly.CSharp.variableDB_.getName(this.getFieldValue('VAR'), Blockly.VARIABLE_CATEGORY_NAME)
      const argument0 = Blockly.CSharp.valueToCode(this, 'FROM', Blockly.CSharp.ORDER_ASSIGNMENT) || '0'
      const argument1 = Blockly.CSharp.valueToCode(this, 'TO', Blockly.CSharp.ORDER_ASSIGNMENT) || '0'
      const increment = Blockly.CSharp.valueToCode(this, 'BY', Blockly.CSharp.ORDER_ASSIGNMENT) || '1'
      let branch = Blockly.CSharp.statementToCode(this, 'DO')
      if (Blockly.CSharp.INFINITE_LOOP_TRAP) {
        branch = Blockly.CSharp.INFINITE_LOOP_TRAP.replace(/%1/g, '\'' + this.id + '\'') + branch
      }
      let code
      if (Blockly.isNumber(argument0) && Blockly.isNumber(argument1) && Blockly.isNumber(increment)) {
        // All arguments are simple numbers.
        const up = parseFloat(argument0) <= parseFloat(argument1)
        code = 'for (' + variable0 + ' = ' + argument0 + '; ' + variable0 + (up ? ' <= ' : ' >= ') + argument1 + '; ' + variable0
        const step = Math.abs(parseFloat(increment))
        if (step === 1) {
          code += up ? '++' : '--'
        } else {
          code += (up ? ' += ' : ' -= ') + step
        }
        code += ') {\n' + branch + '}\n'
      } else {
        code = ''
        // Cache non-trivial values to variables to prevent repeated look-ups.
        let startVar = argument0
        if (!argument0.match(/^\w+$/) && !Blockly.isNumber(argument0)) {
          startVar = Blockly.CSharp.variableDB_.getDistinctName(variable0 + '_start', Blockly.VARIABLE_CATEGORY_NAME)
          code += 'var ' + startVar + ' = ' + argument0 + ';\n'
        }
        let endVar = argument1
        if (!argument1.match(/^\w+$/) && !Blockly.isNumber(argument1)) {
          endVar = Blockly.CSharp.variableDB_.getDistinctName(variable0 + '_end', Blockly.VARIABLE_CATEGORY_NAME)
          code += 'var ' + endVar + ' = ' + argument1 + ';\n'
        }
        // Determine loop direction at start, in case one of the bounds
        // changes during loop execution.
        const incVar = Blockly.CSharp.variableDB_.getDistinctName(variable0 + '_inc', Blockly.VARIABLE_CATEGORY_NAME)
        code += 'var ' + incVar + ' = '
        if (Blockly.isNumber(increment)) {
          code += Math.abs(increment) + ';\n'
        } else {
          code += 'Math.Abs(' + increment + ');\n'
        }
        code += 'if (' + startVar + ' > ' + endVar + ') {\n'
        code += '  ' + incVar + ' = -' + incVar + ';\n'
        code += '}\n'
        code += 'for (' + variable0 + ' = ' + startVar + ';\n' + '     ' + incVar + ' >= 0 ? ' + variable0 + ' <= ' + endVar + ' : ' +
          variable0 + ' >= ' + endVar + ';\n' + '     ' + variable0 + ' += ' + incVar + ') {\n' + branch + '}\n'
      }
      return code
    }

    Blockly.CSharp.controls_forEach = function () {
      // For each loop.
      const variable0 = Blockly.CSharp.variableDB_.getName(this.getFieldValue('VAR'), Blockly.VARIABLE_CATEGORY_NAME)
      const argument0 = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_ASSIGNMENT) || '[]'
      let branch = Blockly.CSharp.statementToCode(this, 'DO')
      if (Blockly.CSharp.INFINITE_LOOP_TRAP) {
        branch = Blockly.CSharp.INFINITE_LOOP_TRAP.replace(/%1/g, '\'' + this.id + '\'') + branch
      }
      let code
      if (argument0.match(/^\w+$/)) {
        code = 'foreach (var ' + variable0 + ' in  ' + argument0 + ') {\n' + branch + '}\n'
      } else {
        // The list appears to be more complicated than a simple variable.
        // Cache it to a variable to prevent repeated look-ups.
        const listVar = Blockly.CSharp.variableDB_.getDistinctName(variable0 + '_list', Blockly.VARIABLE_CATEGORY_NAME)
        code = 'var ' + listVar + ' = ' + argument0 + ';\n' + 'foreach (var ' + variable0 + ' in ' + listVar + ') {\n' + branch + '}\n'
      }
      return code
    }

    Blockly.CSharp.controls_flow_statements = function () {
      // Flow statements: continue, break.
      switch (this.getFieldValue('FLOW')) {
        case 'BREAK':
          return 'break;\n'
        case 'CONTINUE':
          return 'continue;\n'
      }
      throw Error('Unknown flow statement.')
    }
  },
  initMath: function (Blockly) {
    Blockly.CSharp.math = {}

    Blockly.CSharp.math_number = function () {
      // Numeric value.
      var code = window.parseFloat(this.getFieldValue('NUM'))
      if (code.toString().indexOf('.') !== -1) {
        code += 'f'
      }
      return [code, Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp.math_arithmetic = function () {
      // Basic arithmetic operators, and power.
      const mode = this.getFieldValue('OP')
      const tuple = Blockly.CSharp.math_arithmetic.OPERATORS[mode]
      const operator = tuple[0]
      const order = tuple[1]
      const argument0 = Blockly.CSharp.valueToCode(this, 'A', order) || '0.0'
      const argument1 = Blockly.CSharp.valueToCode(this, 'B', order) || '0.0'
      let code
      // Power in cSharp requires a special case since it has no operator.
      if (!operator) {
        code = 'Math.Pow(' + argument0 + ', ' + argument1 + ')'
        return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      }
      code = argument0 + operator + argument1
      return [code, order]
    }

    Blockly.CSharp.math_arithmetic.OPERATORS = {
      ADD: [' + ', Blockly.CSharp.ORDER_ADDITION],
      MINUS: [' - ', Blockly.CSharp.ORDER_SUBTRACTION],
      MULTIPLY: [' * ', Blockly.CSharp.ORDER_MULTIPLICATION],
      DIVIDE: [' / ', Blockly.CSharp.ORDER_DIVISION],
      POWER: [null, Blockly.CSharp.ORDER_COMMA] // Handle power separately.
    }

    Blockly.CSharp.math_single = function () {
      // Math operators with single operand.
      const operator = this.getFieldValue('OP')
      let code
      let arg
      if (operator === 'NEG') {
        // Negation is a special case given its different operator precedence.
        arg = Blockly.CSharp.valueToCode(this, 'NUM', Blockly.CSharp.ORDER_UNARY_NEGATION) || '0.0'
        if (arg[0] === '-') {
          // --3 is not allowed
          arg = ' ' + arg
        }
        code = '-' + arg
        return [code, Blockly.CSharp.ORDER_UNARY_NEGATION]
      }
      if (operator === 'SIN' || operator === 'COS' || operator === 'TAN') {
        arg = Blockly.CSharp.valueToCode(this, 'NUM', Blockly.CSharp.ORDER_DIVISION) || '0'
      } else {
        arg = Blockly.CSharp.valueToCode(this, 'NUM', Blockly.CSharp.ORDER_NONE) || '0.0'
      }
      // First, handle cases which generate values that don't need parentheses
      // wrapping the code.
      switch (operator) {
        case 'ABS':
          code = 'Math.Abs(' + arg + ')'
          break
        case 'ROOT':
          code = 'Math.Sqrt(' + arg + ')'
          break
        case 'LN':
          code = 'Math.Log(' + arg + ')'
          break
        case 'LOG10':
          code = 'Math.Log10(' + arg + ')'
          break
        case 'EXP':
          code = 'Math.Exp(' + arg + ')'
          break
        case 'POW10':
          code = 'Math.Pow(' + arg + ', 10)'
          break
        case 'ROUND':
          code = 'Math.Round(' + arg + ')'
          break
        case 'ROUNDUP':
          code = 'Math.Ceil(' + arg + ')'
          break
        case 'ROUNDDOWN':
          code = 'Math.Floor(' + arg + ')'
          break
        case 'SIN':
          code = 'Math.Sin(' + arg + ' / 180 * Math.PI)'
          break
        case 'COS':
          code = 'Math.Cos(' + arg + ' / 180 * Math.PI)'
          break
        case 'TAN':
          code = 'Math.Tan(' + arg + ' / 180 * Math.PI)'
          break
      }
      if (code) {
        return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      }
      // Second, handle cases which generate values that may need parentheses
      // wrapping the code.
      switch (operator) {
        case 'ASIN':
          code = 'Math.Asin(' + arg + ') / Math.PI * 180'
          break
        case 'ACOS':
          code = 'Math.Acos(' + arg + ') / Math.PI * 180'
          break
        case 'ATAN':
          code = 'Math.Atan(' + arg + ') / Math.PI * 180'
          break
        default:
          throw Error('Unknown math operator: ' + operator)
      }
      return [code, Blockly.CSharp.ORDER_DIVISION]
    }

    Blockly.CSharp.math_constant = function () {
      // Constants: PI, E, the Golden Ratio, sqrt(2), 1/sqrt(2), INFINITY.
      const constant = this.getFieldValue('CONSTANT')
      return Blockly.CSharp.math_constant.CONSTANTS[constant]
    }

    Blockly.CSharp.math_constant.CONSTANTS = {
      PI: ['Math.PI', Blockly.CSharp.ORDER_MEMBER],
      E: ['Math.E', Blockly.CSharp.ORDER_MEMBER],
      GOLDEN_RATIO: ['(1 + Math.Sqrt(5)) / 2', Blockly.CSharp.ORDER_DIVISION],
      SQRT2: ['Math.Sqrt(2)', Blockly.CSharp.ORDER_MEMBER],
      SQRT1_2: ['Math.Sqrt(1.0 / 2)', Blockly.CSharp.ORDER_MEMBER],
      INFINITY: ['double.PositiveInfinity', Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp.math_number_property = function () {
      // Check if a number is even, odd, prime, whole, positive, or negative
      // or if it is divisible by certain number. Returns true or false.
      const numberToCheck = Blockly.CSharp.valueToCode(this, 'NUMBER_TO_CHECK', Blockly.CSharp.ORDER_MODULUS) || 'double.NaN'
      const dropdownProperty = this.getFieldValue('PROPERTY')
      let code
      if (dropdownProperty === 'PRIME') {
        // Prime is a special case as it is not a one-liner test.
        if (!Blockly.CSharp.definitions_['isPrime']) {
          const functionName = Blockly.CSharp.variableDB_.getDistinctName('isPrime', Blockly.Generator.NAME_TYPE)
          Blockly.CSharp.logic_prime = functionName
          const func = []
          func.push('var ' + functionName + ' = new Func<double, bool>((n) => {')
          func.push('  // http://en.wikipedia.org/wiki/Primality_test#Naive_methods')
          func.push('  if (n == 2.0 || n == 3.0)')
          func.push('    return true;')
          func.push('  // False if n is NaN, negative, is 1, or not whole. And false if n is divisible by 2 or 3.')
          func.push('  if (double.IsNaN(n) || n <= 1 || n % 1 != 0.0 || n % 2 == 0.0 || n % 3 == 0.0)')
          func.push('    return false;')
          func.push('  // Check all the numbers of form 6k +/- 1, up to sqrt(n).')
          func.push('  for (var x = 6; x <= Math.Sqrt(n) + 1; x += 6) {')
          func.push('    if (n % (x - 1) == 0.0 || n % (x + 1) == 0.0)')
          func.push('      return false;')
          func.push('  }')
          func.push('  return true;')
          func.push('});')
          Blockly.CSharp.definitions_['isPrime'] = func.join('\n')
        }
        code = Blockly.CSharp.logic_prime + '(' + numberToCheck + ')'
        return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      }
      switch (dropdownProperty) {
        case 'EVEN':
          code = numberToCheck + ' % 2 == 0'
          break
        case 'ODD':
          code = numberToCheck + ' % 2 == 1'
          break
        case 'WHOLE':
          code = numberToCheck + ' % 1 == 0'
          break
        case 'POSITIVE':
          code = numberToCheck + ' > 0'
          break
        case 'NEGATIVE':
          code = numberToCheck + ' < 0'
          break
        case 'DIVISIBLE_BY':
          const divisor = Blockly.CSharp.valueToCode(this, 'DIVISOR', Blockly.CSharp.ORDER_MODULUS) || 'double.NaN'
          code = numberToCheck + ' % ' + divisor + ' == 0'
          break
      }
      return [code, Blockly.CSharp.ORDER_EQUALITY]
    }

    Blockly.CSharp.math_change = function () {
      // Add to a variable in place.
      const argument0 = Blockly.CSharp.valueToCode(this, 'DELTA', Blockly.CSharp.ORDER_ADDITION) || '0.0'
      const varName = Blockly.CSharp.variableDB_.getName(this.getFieldValue('VAR'), Blockly.VARIABLE_CATEGORY_NAME)
      return varName + ' = (' + varName + '.GetType().Name == "Double" ? ' + varName + ' : 0.0) + ' + argument0 + ';\n'
    }

    // Rounding functions have a single operand.
    Blockly.CSharp.math_round = Blockly.CSharp.math_single
    // Trigonometry functions have a single operand.
    Blockly.CSharp.math_trig = Blockly.CSharp.math_single

    Blockly.CSharp.math_on_list = function () {
      // Math functions for lists.
      const func = this.getFieldValue('OP')
      let list, code
      switch (func) {
        case 'SUM':
          list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_MEMBER) || 'new List<dynamic>()'
          code = list + '.Aggregate((x, y) => x + y)'
          break
        case 'MIN':
          list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_COMMA) || 'new List<dynamic>()'
          code = list + '.Min()'
          break
        case 'MAX':
          list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_COMMA) || 'new List<dynamic>()'
          code = list + '.Max()'
          break
        case 'AVERAGE':
          list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_COMMA) || 'new List<dynamic>()'
          code = list + '.Average()'
          break
        case 'MEDIAN':
          // math_median([null,null,1,3]) == 2.0.
          if (!Blockly.CSharp.definitions_['math_median']) {
            const functionName = Blockly.CSharp.variableDB_.getDistinctName('math_median', Blockly.Generator.NAME_TYPE)
            Blockly.CSharp.math_on_list.math_median = functionName
            const func = []
            func.push('var ' + functionName + ' = new Func<List<dynamic>,dynamic>((vals) => {')
            func.push('  vals.Sort();')
            func.push('  if (vals.Count % 2 == 0)')
            func.push('    return (vals[vals.Count / 2 - 1] + vals[vals.Count / 2]) / 2;')
            func.push('  else')
            func.push('    return vals[(vals.Count - 1) / 2];')
            func.push('});')
            Blockly.CSharp.definitions_['math_median'] = func.join('\n')
          }
          list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_NONE) || 'new List<dynamic>()'
          code = Blockly.CSharp.math_on_list.math_median + '(' + list + ')'
          break
        case 'MODE':
          if (!Blockly.CSharp.definitions_['math_modes']) {
            const functionName = Blockly.CSharp.variableDB_.getDistinctName('math_modes', Blockly.Generator.NAME_TYPE)
            Blockly.CSharp.math_on_list.math_modes = functionName
            // As a list of numbers can contain more than one mode,
            // the returned result is provided as an array.
            // Mode of [3, 'x', 'x', 1, 1, 2, '3'] -> ['x', 1].
            const func = []
            func.push('var ' + functionName + ' = new Func<List<dynamic>,List<dynamic>>((values) => {')
            func.push('  var modes = new List<dynamic>();')
            func.push('  var counts = new Dictionary<double, int>();')
            func.push('  var maxCount = 0;')
            func.push('  foreach (var value in values) {')
            func.push('    int storedCount;')
            func.push('    if (counts.TryGetValue(value, out storedCount)) {')
            func.push('      maxCount = Math.Max(maxCount, ++counts[value]);')
            func.push('    }')
            func.push('    else {')
            func.push('      counts.Add(value, 1);')
            func.push('      maxCount = 1;')
            func.push('    }')
            func.push('  }')
            func.push('  foreach (var pair in counts) {')
            func.push('    if (pair.Value == maxCount)')
            func.push('      modes.Add(pair.Key);')
            func.push('  }')
            func.push('  return modes;')
            func.push('});')
            Blockly.CSharp.definitions_['math_modes'] = func.join('\n')
          }
          list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_NONE) || 'new List<dynamic>()'
          code = Blockly.CSharp.math_on_list.math_modes + '(' + list + ')'
          break
        case 'STD_DEV':
          if (!Blockly.CSharp.definitions_['math_standard_deviation']) {
            const functionName = Blockly.CSharp.variableDB_.getDistinctName('math_standard_deviation', Blockly.Generator.NAME_TYPE)
            Blockly.CSharp.math_on_list.math_standard_deviation = functionName
            const func = []
            func.push('var ' + functionName + ' = new Func<List<dynamic>,double>((numbers) => {')
            func.push('  var n = numbers.Count;')
            func.push('  var mean = numbers.Average(val => val);')
            func.push('  var variance = 0.0;')
            func.push('  for (var j = 0; j < n; j++) {')
            func.push('    variance += Math.Pow(numbers[j] - mean, 2);')
            func.push('  }')
            func.push('  variance = variance / n;')
            func.push('  return Math.Sqrt(variance);')
            func.push('});')
            Blockly.CSharp.definitions_['math_standard_deviation'] = func.join('\n')
          }
          list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_NONE) || 'new List<dynamic>()'
          code = Blockly.CSharp.math_on_list.math_standard_deviation + '(' + list + ')'
          break
        case 'RANDOM':
          if (!Blockly.CSharp.definitions_['math_random_item']) {
            const functionName = Blockly.CSharp.variableDB_.getDistinctName('math_random_item', Blockly.Generator.NAME_TYPE)
            Blockly.CSharp.math_on_list.math_random_item = functionName
            const func = []
            func.push('var ' + functionName + ' = new Func<List<dynamic>,dynamic>((list) => {')
            func.push('  var x = (new Random()).Next(list.Count);')
            func.push('  return list[x];')
            func.push('});')
            Blockly.CSharp.definitions_['math_random_item'] = func.join('\n')
          }
          list = Blockly.CSharp.valueToCode(this, 'LIST', Blockly.CSharp.ORDER_NONE) || 'new List<dynamic>()'
          code = Blockly.CSharp.math_on_list.math_random_item + '(' + list + ')'
          break
        default:
          throw Error('Unknown operator: ' + func)
      }
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.math_modulo = function () {
      // Remainder computation.
      const argument0 = Blockly.CSharp.valueToCode(this, 'DIVIDEND', Blockly.CSharp.ORDER_MODULUS) || '0.0'
      const argument1 = Blockly.CSharp.valueToCode(this, 'DIVISOR', Blockly.CSharp.ORDER_MODULUS) || '0.0'
      const code = argument0 + ' % ' + argument1
      return [code, Blockly.CSharp.ORDER_MODULUS]
    }

    Blockly.CSharp.math_constrain = function () {
      // Constrain a number between two limits.
      const argument0 = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_COMMA) || '0.0'
      const argument1 = Blockly.CSharp.valueToCode(this, 'LOW', Blockly.CSharp.ORDER_COMMA) || '0.0'
      const argument2 = Blockly.CSharp.valueToCode(this, 'HIGH', Blockly.CSharp.ORDER_COMMA) || 'double.PositiveInfinity'
      const code = 'Math.Min(Math.Max(' + argument0 + ', ' + argument1 + '), ' + argument2 + ')'
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.math_random_int = function () {
      // Random integer between [X] and [Y].
      const argument0 = Blockly.CSharp.valueToCode(this, 'FROM', Blockly.CSharp.ORDER_COMMA) || '0.0'
      const argument1 = Blockly.CSharp.valueToCode(this, 'TO', Blockly.CSharp.ORDER_COMMA) || '0.0'

      const code = 'Utils.RandomInt(' + argument0 + ', ' + argument1 + ')'
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.math_random_float = function () {
      // Random fraction between 0 and 1.
      return ['Utils.RandomDouble()', Blockly.CSharp.ORDER_FUNCTION_CALL]
    }
  },
  initProcedures: function (Blockly) {
    Blockly.CSharp.procedures = {}

    Blockly.CSharp.procedures_defreturn = function () {
      // Define a procedure with a return value.
      const funcName = Blockly.CSharp.variableDB_.getName(this.getFieldValue('NAME'), Blockly.PROCEDURE_CATEGORY_NAME)
      let branch = Blockly.CSharp.statementToCode(this, 'STACK')

      if (Blockly.CSharp.INFINITE_LOOP_TRAP) {
        branch = Blockly.CSharp.INFINITE_LOOP_TRAP.replace(/%1/g, '\'' + this.id + '\'') + branch
      }

      let returnValue = Blockly.CSharp.valueToCode(this, 'RETURN', Blockly.CSharp.ORDER_NONE) || ''
      if (returnValue) {
        returnValue = '  return ' + returnValue + ';\n'
      }

      let args = []
      for (let x = 0; x < this.arguments_.length; x++) {
        args[x] = Blockly.CSharp.variableDB_.getName(this.arguments_[x], Blockly.VARIABLE_CATEGORY_NAME)
      }

      let argTypes = ''
      const appendToList = function (res, val) {
        if (res.length === 0) {
          argTypes = val
        } else {
          argTypes += ', ' + val
        }
      }

      for (let x = 0; x < args.length; x++) {
        appendToList(argTypes, 'dynamic')
      }

      if (returnValue.length !== 0) {
        appendToList(argTypes, 'dynamic')
      }

      const delegateType = (returnValue.length === 0) ? 'Action' : ('Func<' + argTypes + '>')

      let code = 'var ' + funcName + ' = new ' + delegateType + '((' + args.join(', ') + ') => {\n' + branch + returnValue + '});'
      code = Blockly.CSharp.scrub_(this, code)
      Blockly.CSharp.definitions_[funcName] = code
      return null
    }

    // Defining a procedure without a return value uses the same generator as
    // a procedure with a return value.
    Blockly.CSharp.procedures_defnoreturn = Blockly.CSharp.procedures_defreturn

    Blockly.CSharp.procedures_callreturn = function () {
      // Call a procedure with a return value.
      const funcName = Blockly.CSharp.variableDB_.getName(this.getFieldValue('NAME'), Blockly.PROCEDURE_CATEGORY_NAME)
      let args = []
      for (let x = 0; x < this.arguments_.length; x++) {
        args[x] = Blockly.CSharp.valueToCode(this, 'ARG' + x, Blockly.CSharp.ORDER_COMMA) || 'null'
      }
      const code = funcName + '(' + args.join(', ') + ')'
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.procedures_callnoreturn = function () {
      // Call a procedure with no return value.
      const funcName = Blockly.CSharp.variableDB_.getName(this.getFieldValue('NAME'), Blockly.PROCEDURE_CATEGORY_NAME)
      let args = []
      for (let x = 0; x < this.arguments_.length; x++) {
        args[x] = Blockly.CSharp.valueToCode(this, 'ARG' + x, Blockly.CSharp.ORDER_COMMA) || 'null'
      }
      return funcName + '(' + args.join(', ') + ');\n'
    }

    Blockly.CSharp.procedures_ifreturn = function () {
      // Conditionally return value from a procedure.
      const condition = Blockly.CSharp.valueToCode(this, 'CONDITION', Blockly.CSharp.ORDER_NONE) || 'false'
      let code = 'if (' + condition + ') {\n'
      if (this.hasReturnValue_) {
        const value = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_NONE) || 'null'
        code += '  return ' + value + ';\n'
      } else {
        code += '  return;\n'
      }
      code += '}\n'
      return code
    }
  },
  initText: function (Blockly) {
    Blockly.CSharp.text = {}

    Blockly.CSharp.text = function () {
      // Text value.
      const code = Blockly.CSharp.quote_(this.getFieldValue('TEXT'))
      return [code, Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp.text_join = function () {
      if (this.itemCount_ === 0) {
        return ['""', Blockly.CSharp.ORDER_ATOMIC]
      } else if (this.itemCount_ === 1) {
        const argument0 = Blockly.CSharp.valueToCode(this, 'ADD0', Blockly.CSharp.ORDER_NONE) || '""'
        const code = argument0 + '.ToString()'
        return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      } else if (this.itemCount_ === 2) {
        const argument0 = Blockly.CSharp.valueToCode(this, 'ADD0', Blockly.CSharp.ORDER_NONE) || '""'
        const argument1 = Blockly.CSharp.valueToCode(this, 'ADD1', Blockly.CSharp.ORDER_NONE) || '""'
        const code = 'String.Concat(' + argument0 + ', ' + argument1 + ')'
        return [code, Blockly.CSharp.ORDER_ADDITION]
      } else {
        let code = new Array(this.itemCount_)
        for (let n = 0; n < this.itemCount_; n++) {
          code[n] = Blockly.CSharp.valueToCode(this, 'ADD' + n, Blockly.CSharp.ORDER_COMMA) || '""'
        }
        code = 'String.Concat(' + code.join(', ') + ')'
        return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      }
    }

    Blockly.CSharp.text_append = function () {
      // Append to a variable in place.
      const varName = Blockly.CSharp.variableDB_.getName(this.getFieldValue('VAR'), Blockly.VARIABLE_CATEGORY_NAME)
      const argument0 = Blockly.CSharp.valueToCode(this, 'TEXT', Blockly.CSharp.ORDER_NONE) || '""'
      return varName + ' = String.Concat(' + varName + ', ' + argument0 + ');\n'
    }

    Blockly.CSharp.text_length = function () {
      // String length.
      const argument0 = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_FUNCTION_CALL) || '""'
      return [argument0 + '.Length', Blockly.CSharp.ORDER_MEMBER]
    }

    Blockly.CSharp.text_isEmpty = function () {
      // Is the string null?
      const argument0 = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_MEMBER) || '""'
      return [argument0 + '.Length == 0', Blockly.CSharp.ORDER_EQUALITY]
    }

    Blockly.CSharp.text_indexOf = function () {
      // Search the text for a substring.
      const operator = this.getFieldValue('END') === 'FIRST' ? 'IndexOf' : 'LastIndexOf'
      const argument0 = Blockly.CSharp.valueToCode(this, 'FIND', Blockly.CSharp.ORDER_NONE) || '""'
      const argument1 = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_MEMBER) || '""'
      const code = argument1 + '.' + operator + '(' + argument0 + ') + 1'
      return [code, Blockly.CSharp.ORDER_MEMBER]
    }

    Blockly.CSharp.text_charAt = function () {
      const where = this.getFieldValue('WHERE') || 'FROM_START'
      const at = Blockly.CSharp.valueToCode(this, 'AT', Blockly.CSharp.ORDER_UNARY_NEGATION) || '1'
      const text = Blockly.CSharp.valueToCode(this, 'VALUE', Blockly.CSharp.ORDER_MEMBER) || '""'

      // Blockly uses one-based indicies.
      let code
      switch (where) {
        case 'FIRST':
          code = text + '.First()'
          return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
        case 'LAST':
          code = text + '.Last()'
          return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
        case 'FROM_START':
          code = text + '[' + at + ' - 1]'
          return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
        case 'FROM_END':
          code = text + '[text.Length - ' + at + ']'
          return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
        case 'RANDOM':
          if (!Blockly.CSharp.definitions_['text_random_letter']) {
            const functionName = Blockly.CSharp.variableDB_.getDistinctName('text_random_letter', Blockly.Generator.NAME_TYPE)
            Blockly.CSharp.text_charAt.text_random_letter = functionName
            const func = []
            func.push('var ' + functionName + ' = new Func<string, char>((text) => {')
            func.push('  var x = (new Random()).Next(text.length);')
            func.push('  return text[x];')
            func.push('});')
            Blockly.CSharp.definitions_['text_random_letter'] = func.join('\n')
          }
          code = Blockly.CSharp.text_charAt.text_random_letter + '(' + text + ')'
          return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
      }
      throw Error('Unhandled option (text_charAt).')
    }

    Blockly.CSharp.text_getSubstring = function () {
      // Get substring.
      const text = Blockly.CSharp.valueToCode(this, 'STRING', Blockly.CSharp.ORDER_MEMBER) || 'null'
      const where1 = this.getFieldValue('WHERE1')
      const where2 = this.getFieldValue('WHERE2')
      const at1 = Blockly.CSharp.valueToCode(this, 'AT1', Blockly.CSharp.ORDER_NONE) || '1'
      const at2 = Blockly.CSharp.valueToCode(this, 'AT2', Blockly.CSharp.ORDER_NONE) || '1'
      let code
      if (where1 === 'FIRST' && where2 === 'LAST') {
        code = text
      } else {
        if (!Blockly.CSharp.definitions_['text_get_substring']) {
          const functionName = Blockly.CSharp.variableDB_.getDistinctName('text_get_substring', Blockly.Generator.NAME_TYPE)
          Blockly.CSharp.text_getSubstring.func = functionName
          const func = []
          func.push('var ' + functionName + ' = new Func<string, dynamic, int, dynamic, int, string>((text, where1, at1, where2, at2) => {')
          func.push('  var getAt =new Func<dynamic, int, int>((where, at) => {')
          func.push('    if (where == "FROM_START") {')
          func.push('      at--;')
          func.push('    } else if (where == "FROM_END") {')
          func.push('      at = text.Length - at;')
          func.push('    } else if (where == "FIRST") {')
          func.push('      at = 0;')
          func.push('    } else if (where == "LAST") {')
          func.push('      at = text.Length - 1;')
          func.push('    } else {')
          func.push('      throw new ApplicationException("Unhandled option (text_getSubstring).");')
          func.push('    }')
          func.push('    return at;')
          func.push('  });')
          func.push('  at1 = getAt(where1, at1);')
          func.push('  at2 = getAt(where2, at2) + 1;')
          func.push('  return text.Substring(at1, at2 - at1);')
          func.push('});')
          Blockly.CSharp.definitions_['text_get_substring'] = func.join('\n')
        }
        code = Blockly.CSharp.text_getSubstring.func + '(' + text + ', "' + where1 + '", ' + at1 + ', "' + where2 + '", ' + at2 + ')'
      }
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.text_changeCase = function () {
      // Change capitalization.
      const mode = this.getFieldValue('CASE')
      const operator = Blockly.CSharp.text_changeCase.OPERATORS[mode]
      let code
      if (operator) {
        // Upper and lower case are functions built into cSharp.
        const argument0 = Blockly.CSharp.valueToCode(this, 'TEXT', Blockly.CSharp.ORDER_MEMBER) || '""'
        code = argument0 + operator
      } else {
        if (!Blockly.CSharp.definitions_['text_toTitleCase']) {
          // Title case is not a native cSharp function.  Define one.
          const functionName = Blockly.CSharp.variableDB_.getDistinctName('text_toTitleCase', Blockly.Generator.NAME_TYPE)
          Blockly.CSharp.text_changeCase.toTitleCase = functionName
          const func = []
          func.push('var ' + functionName + ' = new Func<string, string>((str) => {')
          func.push('  var buf = new StringBuilder(str.Length);')
          func.push('  var toUpper = true;')
          func.push('  foreach (var ch in str) {')
          func.push('    buf.Append(toUpper ? Char.ToUpper(ch) : ch);')
          func.push('    toUpper = Char.IsWhiteSpace(ch);')
          func.push('  }')
          func.push('  return buf.ToString();')
          func.push('});')
          Blockly.CSharp.definitions_['text_toTitleCase'] = func.join('\n')
        }
        const argument0 = Blockly.CSharp.valueToCode(this, 'TEXT', Blockly.CSharp.ORDER_NONE) || '""'
        code = Blockly.CSharp.text_changeCase.toTitleCase + '(' + argument0 + ')'
      }
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.text_changeCase.OPERATORS = {
      UPPERCASE: '.ToUpper()',
      LOWERCASE: '.ToLower()',
      TITLECASE: null
    }

    Blockly.CSharp.text_trim = function () {
      // Trim spaces.
      const mode = this.getFieldValue('MODE')
      const operator = Blockly.CSharp.text_trim.OPERATORS[mode]
      const argument0 = Blockly.CSharp.valueToCode(this, 'TEXT', Blockly.CSharp.ORDER_MEMBER) || '""'
      return [argument0 + operator, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }

    Blockly.CSharp.text_trim.OPERATORS = {
      LEFT: '.TrimStart()',
      RIGHT: '.TrimEnd()',
      BOTH: '.Trim()'
    }

    Blockly.CSharp.text_print = function () {
      // Print statement.
      const argument0 = Blockly.CSharp.valueToCode(this, 'TEXT', Blockly.CSharp.ORDER_NONE) || '""'
      return 'Console.WriteLine(' + argument0 + ');\n'
    }

    Blockly.CSharp.text_prompt = function () {
      const msg = Blockly.CSharp.quote_(this.getFieldValue('TEXT'))
      const toNumber = this.getFieldValue('TYPE') === 'NUMBER'

      const functionName = Blockly.CSharp.variableDB_.getDistinctName('text_promptInput', Blockly.Generator.NAME_TYPE)
      Blockly.CSharp.text_prompt.promptInput = functionName
      const func = []
      func.push('var ' + functionName + ' = new Func<string, bool, dynamic>((msg, toNumber) => {')
      func.push('  Console.WriteLine(msg);')
      func.push('  var res = Console.ReadLine();')
      func.push('  if (toNumber)')
      func.push('    return Double.Parse(res);')
      func.push('  return res;')
      func.push('});')
      Blockly.CSharp.definitions_['text_promptInput'] = func.join('\n')

      const code = Blockly.CSharp.text_prompt.promptInput + '(' + msg + ', ' + toNumber + ')'
      return [code, Blockly.CSharp.ORDER_FUNCTION_CALL]
    }
  },
  initVariables: function (Blockly) {
    Blockly.CSharp.variables = {}

    // Override variables
    Blockly.CSharp['variables_get'] = Blockly.CSharp['variables_get_dynamic'] = function (block) {
      // Variable getter.
      const code = Blockly.CSharp.variableDB_.getName(block.getFieldValue('VAR'), Blockly.VARIABLE_CATEGORY_NAME)
      return [code, Blockly.CSharp.ORDER_ATOMIC]
    }

    Blockly.CSharp['variables_set'] = Blockly.CSharp['variables_set_dynamic'] = function (block) {
      // Variable setter.
      const argument0 = Blockly.CSharp.valueToCode(block, 'VALUE', Blockly.CSharp.ORDER_ASSIGNMENT) || '0'
      const varName = Blockly.CSharp.variableDB_.getName(block.getFieldValue('VAR'), Blockly.VARIABLE_CATEGORY_NAME)
      return varName + ' = ' + argument0 + ';\n'
    }
  }
}
