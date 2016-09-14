source ./go/gocompiler.vim
source ./go/typedefs.vim

call extend(s:, ImportGoCompiler())

function! s:generate()
  let gofile = 'go/vimlparser.go'
  " let vimlfunc = 'go/vimlfunc/vimlfunc.go'
  let head = [
  \   "// Code generated by gocompiler.vim",
  \   "// source: autoload/vimlparser.vim",
  \   "// DO NOT EDIT!",
  \   "",
  \   "package vimlparser",
  \ ]

  try
    let ast = s:ast()
    let c = s:GoCompiler.new(ImportTypedefs())
    let lines = c.compile(ast)
    call writefile(head + lines, gofile)
  catch
    echoerr substitute(v:throwpoint, '\.\.\zs\d\+', '\=s:numtoname(submatch(0))', 'g') . "\n" . v:exception
  endtry
endfunction

function! s:ast() abort
  let vimfile = 'autoload/vimlparser.vim'
  let astfile = 'go/vimlparser.ast.vim'

  let lines = readfile(vimfile)
  unlet lines[0:index(lines, 'let s:FALSE = 0')]
  unlet lines[index(lines, 'let s:RegexpParser = {}'):-2]
  let r = s:StringReader.new(lines)
  let p = s:VimLParser.new()
  let ast = p.parse(r)
  return ast
endfunction

function! s:numtoname(num)
  let sig = printf("function('%s')", a:num)
  for k in keys(s:)
    if type(s:[k]) == type({})
      for name in keys(s:[k])
        if type(s:[k][name]) == type(function('tr')) && string(s:[k][name]) == sig
          return printf('%s.%s', k, name)
        endif
      endfor
    endif
  endfor
  return a:num
endfunction

call s:generate()
