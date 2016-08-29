tinymce.init({
    selector:'textarea',
    setup: function (ed) {
        ed.on('init', function(args) {
            var id = ed.id;
            var height = 300;

            document.getElementById(id + '_ifr').style.height = height + 'px';
            document.getElementById(id + '_tbl').style.height = (height + 30) + 'px';
        });
    }});