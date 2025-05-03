import Big from 'big.js';
import 'bulma/css/bulma.css';
import $ from 'jquery';
// make jQuery global (for esbuild)
window.$ = $;
window.jQuery = $;

// recommended way of calling $(document).ready();
$(function() {
    console.log("DOM ready with jquery installed.");
});
