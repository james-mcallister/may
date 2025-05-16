import Big from 'big.js';
import 'bulma/css/bulma.css';
import '@fortawesome/fontawesome-free';
import $ from 'jquery';
// make jQuery global (for esbuild)
window.$ = $;
window.jQuery = $;

// recommended way of calling $(document).ready();
$(function() {
    HomeModule.init()
});

const HomeModule = (function($) {
    const ele = {
        navbar: $("nav")
    };

    function onClick(e) {
        e.stopPropagation();
        e.preventDefault();
        // call teardown for the current module
        // call init for the selected module
        // need a data structure to store a mapping between url and js module
        console.log($(this).text());
    }

    return {
        init() {
            ele.navbar.on("click", "a", onClick);
        },
        teardown() {
            ele.navbar.off("click", "a", onClick);
        }
    }
})(jQuery);
