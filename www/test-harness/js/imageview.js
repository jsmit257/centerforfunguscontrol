$(_ => {
  $('body>div.imageview')
    .on('activate', (e, src) => {
      $(e.currentTarget)
        .addClass('active')
        .find('img.viewport')
        .attr('src', src)
    })
    .on('click', e => {
      $(e.currentTarget).removeClass('active')
    })
})