$(function () {

  $('.table>.columns').on('clear-sort', e => {
    $(e.currentTarget).find('[sort-order]').removeAttr('sort-order')
  })

  $('.table>.columns>.column[sort-key]')
    .append($('<div class="sort-icon" />')
      .html('&nbsp'))
    .on('click', e => {
      var $col = $(e.currentTarget)
      var sortOrder = $col.attr('sort-order') !== 'sort-asc' ? 'sort-asc' : 'sort-desc'

      $col.parent().trigger('clear-sort')

      $col
        .attr('sort-order', sortOrder)
        .parents('.table')
        .find('.rows')
        .trigger('sort', [$col.attr('sort-key'), $col.attr('sort-order')])
    })

})
