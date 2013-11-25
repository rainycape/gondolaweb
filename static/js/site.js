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

function setupHashAnchors(el) {
    el = el || $(document);
    el.find('a').each(function () {
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
}

$(function () {
    $("[data-toggle='tooltip']").tooltip(); 
    if (window.location.hash) {
        var el = $(window.location.hash);
        if (el.length) {
            setTimeout(function () {
                scrollTo(el);
            }, 0);
        }
    }
    $('a.slide-up').click(function () {
        var div = $(this).parents('h2').first().next();
        if (div.is(':visible')) {
            div.slideUp();
        } else {
            div.slideDown();
        }
        $(this).toggleClass('inverted');
        return false;
    });
    $('a[rel="popover"]').click(function () {
        var a = $(this);
        if (a.data('popover')) {
            a.popover('destroy');
            a.data('popover', false);
            $('body').off('click.popover');
        } else {
            var target = $(a.attr('href')).clone();
            target.attr('id', null);
            a.popover({
                html: true,
                placement: 'auto',
                content: target.html(),
                trigger: 'manual'
            });
            a.popover('show');
            a.data('popover', true);
            $('body').on('click.popover', function (e) {
                if (a.has(e.target).length === 0 && $('.popover').has(e.target).length === 0) {
                    a.click();
                }
            });
        }
        return false;
    });
});
