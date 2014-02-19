$(function() {
    $('a[href*="#"]').click(function (event) {
      event.preventDefault();
      //calculate destination place
      var dest = 0;
      if ($(this.hash).offset().top > $(document).height() - $(window).height()) {
        dest = $(document).height() - $(window).height();
      } else {
        dest = $(this.hash).offset().top;
      }
      var cur = $(document).scrollTop();
      var ms = Math.max(Math.abs(cur - dest) / 3, 300);
      //go to destination
      $('html,body').animate({
        scrollTop: dest
      }, ms, 'swing');
  });
});
