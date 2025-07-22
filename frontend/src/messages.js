import $ from 'jquery';
window.$ = $;
window.jQuery = $;

let notifyEle;

export function init() {
    notifyEle = $("#notification");
    notifyEle.on("may:notify", function(e, nColor, msg) {
        e.stopPropagation();
        e.preventDefault();
        notifyEle.removeClass("is-hidden");
        let markup = `
        <div class="notification is-${nColor} is-light">${msg}</div>
        `
        
        // 500ms animations with a 5 second delay between reveal and hide.
        notifyEle.html(markup).slideDown(500, function() {
            setTimeout(function() {
            notifyEle.slideUp(500, function() {
                notifyEle.empty();
                notifyEle.addClass("is-hidden");
            });
            }, 5000);
        });
    });
}

export function notify(nColor, msg) {
    $("#notification").trigger("may:notify", [nColor, msg]);
}

export function showProgress() {
    let markup = `<progress class="progress is-small is-link" max="100"></progress>`
    $("#notification").html(markup);
}

export function endProgress() {
    $("#notification").empty();
}
