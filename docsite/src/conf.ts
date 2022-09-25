import hljs from 'highlight.js/lib/core';
import typescript from 'highlight.js/lib/languages/typescript';
import javascript from 'highlight.js/lib/languages/javascript';
import go from 'highlight.js/lib/languages/go';
import python from 'highlight.js/lib/languages/python';
import bash from 'highlight.js/lib/languages/bash';
hljs.registerLanguage('typescript', typescript);
hljs.registerLanguage('javascript', javascript);
hljs.registerLanguage('go', go);
hljs.registerLanguage('bash', bash);
hljs.registerLanguage('python', python);

const libName = "Quid";

const links: Array<{ href: string; name: string }> = [

];

const examplesExtension = ".ts";

export { libName, links, examplesExtension, hljs }