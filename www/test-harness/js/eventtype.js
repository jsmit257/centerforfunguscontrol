$(function () {
  var $eventtype = $('body>.main>.workspace>.eventtype')
  var $table = $eventtype.find('>.table>.rows')
  var $buttonbar = $eventtype.find('>.table>.buttonbar')

  var severities = [
    $('<option value=Begin>Begin</option>'),
    $('<option value=Info>Info</option>'),
    $('<option value=Warn>Warn</option>'),
    $('<option value=Error>Error</option>'),
    $('<option value=Fatal>Fatal</option>'),
    $('<option value=RIP>RIP</option>'),
    $('<option value=Generation>Generation</option>'),
  ]
  var stages = []

  function newRow(data) {
    data ||= { stage: {} }
    return $('<div class="row hover" />')
      .append($('<div class=uuid />').text(data.id))
      .append($('<div class="name static" />').text(data.name))
      .append($('<input class="name live" />').val(data.name))
      .append($('<div class="severity static" />').text(data.severity))
      .append($('<select class="severity live" />')
        .append(severities)
        .val(data.severity))
      .append($('<div class="stage static" />').text(data.stage.name))
      .append($('<select class="stage live" />')
        .append(stages)
        .data('stage_uuid', data.stage.id)
        .val(data.stage.id))
  }

  $table.on('reinit', e => {
    stages = []
    $.ajax({
      url: '/stages',
      method: 'GET',
      async: false,
      success: (result, status, xhr) => {
        result.forEach(r => {
          stages.push($(`<option value="${r.id}">${r.name}</option>`))
        })
      },
      error: console.log,
    })

    $(e.currentTarget).trigger('refresh', {
      newRow: newRow,
      buttonbar: $buttonbar
    })
  })

  $buttonbar.find('>.edit').on('click', e => {
    if (!$(e.currentTarget).hasClass('active')) {
      return
    }
    $table.find('.selected>select.severity')
      .append(severities)
      .val($table.find('.selected>div.severity').text())
    $table.find('.selected>select.stage')
      .append(stages)
      .val($table.find('.selected>select.stage').data('stage_uuid'))

    $table.trigger('edit', {
      data: $selected => {
        return JSON.stringify({
          "name": $selected.find('>.name.live').val(),
          "severity": $selected.find('>.severity.live').val(),
          "stage": {
            "id": $selected.find('>.stage.live').val()
          }
        })
      },
      success: (data, status, xhr) => {
        var $selected = $table.find('.selected')
        $selected.find('>.name.static').text($selected.find('>.name.live').val())
        $selected.find('>.severity.static').text($selected.find('>.severity.live').val())
        $selected
          .find('>.stage.static')
          .text($selected.find('>.stage.live>option:selected').text())
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
          "severity": $selected.find('>.severity.live').val(),
          "stage": {
            "id": $selected.find('>.stage.live').val()
          }
        })
      },
      success: (data, status, xhr) => {
        var $selected = $table.find('.selected')
        $selected.find('>.uuid').text(data.id)
        $selected.find('>.name.static').text($selected.find('>.name.live').val())
        $selected.find('>.severity.static').text($selected.find('>.severity.live').val())
        $selected
          .find('>.stage.static')
          .text($selected.find('>.stage.live>option:selected').text())
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
    $table.trigger('reinit')
  })

  $eventtype.on('activate', e => {
    $eventtype.addClass('active')
    $table.trigger('reinit')
  })
})