function addLineLinks() {
    var lines = $('.gutter .line');
    if (!lines.length) {
        return false;
    }
    lines.each(function () {
        var div = $(this);
        if (div.hasClass('highlighted')) {
            scrollTo(div);
        }
        var text = div.text();
        var id = 'line-' + text;
        var a = $('<a href="#' + id + '" id="' + id + '">' + text + '</a>');
        div.empty();
        div.append(a);
        a.click(function () {
            $('.syntaxhighlighter .highlighted').removeClass('highlighted');
            $('.syntaxhighlighter .number' + $(this).text()).addClass('highlighted');
            var pos = scrollPosition();
            window.location.hash = '#' + $(this).attr('id');
            scrollPosition(pos);
            scrollTo($(this), true);
            return false;
        });
    });
    return true;
}

$(function () {
    var highlight = null;
    var match = /#line\-(\d+)/g.exec(window.location.hash);
    if (match) {
        highlight = parseInt(match[1], 10);
    }
    SyntaxHighlighter.all({toolbar: false, highlight:highlight});
    var ts = setInterval(function () {
        if (addLineLinks()) {
            clearInterval(ts);
        }
    }, 10);
});
