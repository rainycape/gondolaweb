function scrollTo(el, animated) {
    var pos = el.offset().top - 70;
    scrollPosition(pos, animated);
}

function scrollPosition(pos, animated) {
    if (pos !== undefined) {
        if (animated) {
            $('html, body').animate({scrollTop:pos}, 500);
        } else {
            $('html, body').scrollTop(pos);
        }
    }
    return $(document).scrollTop();
}

$(function () {
    if (window.location.hash) {
        var el = $(window.location.hash);
        if (el.length) {
            setTimeout(function () {
                scrollTo(el);
            }, 0);
        }
    }
    $('a').each(function () {
        var a = $(this);
        var href = a.attr('href');
        if (href && href[0] == '#' && href != "#") {
            a.click(function () {
                var target = $($(this).attr('href'));
                if (target.length) {
                    var pos = scrollPosition();
                    window.location.hash = $(this).attr('href');
                    scrollPosition(pos);
                    scrollTo(target, true);
                }
                return false;
            });
        }
    });
});
