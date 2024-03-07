$(function () {
  var $vendor = $('body>.main>.workspace>.vendor')
  var $table = $vendor.find('>.table>.rows')

  $table
    .on('refresh', e => {
      $rows = $(e.currentTarget)
      $.ajax({
        url: '/vendors',
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          var selected = $table.find('.selected>.uuid').text()
          $rows.empty()
          result.forEach(r => {
            var $row = $('<div class="row hover" />')
              .append($('<div class=uuid />').text(r.id))
              .append($('<div class=name />').text(r.name))
              .append($('<div class=website />').text(r.website))
            if (r.id === selected) {
              $row.addClass('selected')
            }
            $rows.append($row)
          })
          if ($rows.find('.selected').length == 0) {
            $rows.find('.row').first().click()
          }
        },
        error: (xhr, status, err) => {
          console.log(xhr, status, err)
        },
      })
    })
    .on('click', '>.row', e => {
      var $row = $(e.currentTarget)
      $row
        .parent()
        .find('.row.selected')
        .removeClass('selected')
      $row.addClass('selected')
    })

  $vendor.on('activate', e => {
    console.log('activating')
    $vendor
      .addClass('active')
      .find('>.table>.rows')
      .trigger('refresh')
  })

})