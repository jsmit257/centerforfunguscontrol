$(function () {
  var $vendor = $('body>.main>.workspace>.vendor')
  var $table = $vendor.find('>.table>.rows')
  var $buttonbar = $vendor.find('>.table>.buttonbar')

  function newRow(data) {
    data ||= {}
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="name static" />').text(data.name))
      .append($('<input class="name live" />').val(data.name))
      .append($('<div class="website static" />')
        .append($('<a />')
          .attr({
            href: data.website,
            target: '_erehwon',
            rel: "noreferrer noopener",
          }).text(data.website)))
      .append($('<input class="website live" />').val(data.website))
  }

  $buttonbar.find('>.edit').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }
    $table.trigger('edit', {
      data: $selected => {
        return JSON.stringify({
          "name": $selected.find('>.name.live').val(),
          "website": $selected.find('>.website.live').val()
        })
      },
      success: (data, status, xhr) => {
        var $selected = $table.find('.selected')
        $selected.find('>.name.static').text($selected.find('>.name.live').val())
        $selected.find('>.website.static>a')
          .attr('href', $selected.find('>.website.live').val())
          .text($selected.find('>.website.live').val())
      },
      buttonbar: $buttonbar
    })
  })

  $buttonbar.find('>.add').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }
    $table.trigger('add', {
      newRow: newRow,
      data: $selected => {
        return JSON.stringify({
          "name": $selected.find('>.name.live').val(),
          "website": $selected.find('>.website.live').val()
        })
      },
      success: (data, status, xhr) => {
        var $selected = $table.find('.selected')
        $selected.find('>.uuid').text(data.id)
        $selected.find('>.name.static').text($selected.find('>.name.live').val())
        $selected.find('>.website.static').text($selected.find('>.website.live').val())
      },
      error: (xhr, status, error) => { $table.trigger('remove-selected') },
      buttonbar: $buttonbar
    })

  })

  $buttonbar.find('>.remove').on('click', e => {
    if ($(e.currentTarget).hasClass('active')) {
      $table.trigger('delete', { buttonbar: $buttonbar })
    }
  })

  $buttonbar.find('>.refresh').on('click', e => {
    $table.trigger('refresh', { newRow: newRow })
  })

  $vendor.on('activate', e => {
    $vendor
      .addClass('active')
      .find('>.table>.rows')
      .trigger('refresh', {
        newRow, newRow,
        buttonbar: $buttonbar
      })
  })
})