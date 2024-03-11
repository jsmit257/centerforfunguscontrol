$(function () {
  var $lifecycle = $('body>.main>.workspace>.lifecycle')
  var fields = ['id', 'location', 'strain_cost', 'grain_cost', 'bulk_cost', 'yield', 'count', 'gross', 'modified_date', 'create_date']

  var $form = $lifecycle.find('form[name=lifecycle]')
    .on('rx', (e, data) => {
      var $items = $(e.target).find('>ul.lifecycle')
      fields.forEach(n => { $items.find(`li.attr[name=${n}]>.input>.static`).text(data[n]) })
      $items
        .find('li.attr[name=strain]>.input>.static')
        .text(`${data.strain.name} | Species: ${data.strain.species || "unknown"} | Vendor: ${data.strain.vendor.name}`)
        .data('strain_uuid', data.strain.id)
      $items
        .find('li.attr[name=grain_substrate]>.input>.static')
        .text(`${data.grain_substrate.name} | Vendor: ${data.grain_substrate.vendor.name}`)
        .data('grainsubstrate_uuid', data.grain_substrate.id)
      $items
        .find('li.attr[name=bulk_substrate]>.input>.static')
        .text(`${data.bulk_substrate.name} | Vendor: ${data.bulk_substrate.vendor.name}`)
        .data('bulksubstrate_uuid', data.bulk_substrate.id)
    })
    .on('tx', (e) => {
      var $items = $(e.target).find('>ul.lifecycle')
      var data = {
        "strain": {
          "id": $items.find('li.attr[name=strain]>.input>.static').data('strain_uuid')
        },
        "grain_substrate": {
          "id": $items.find('li.attr[name=grain_substrate]>.input>.static').data('grainsubstrate_uuid')
        },
        "bulk_substrate": {
          "id": $items.find('li.attr[name=bulk_substrate]>.input>.static').data('bulksubstrate_uuid')
        }
      }
      fields.forEach(n => { data[n] = $items.find(`li.attr[name=${n}]>.input>.live`).val() })
      $.ajax({
        url: "/lifecycle",
        method: 'POST',
        data: data,
        async: true,
        success: (result, status, xhr) => {
        },
        error: (xhr, status, err) => {
        },
      })
    })

  var $ndx = $lifecycle.find('.ndx>.rows')
    .on('refresh', e => {
      var $rows = $(e.currentTarget)
      $.ajax({
        url: '/lifecycles',
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          var selected = $rows.find('.selected>.uuid').text()
          $rows.empty()
          result.forEach(r => {
            var $row = $('<div class="row hover" />')
              .append($('<div class=uuid />').text(r.id))
              .append($('<div class=created_at />').text(r.create_date))
              .append($('<div class=location />').text(r.location))
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
      $.ajax({
        url: '/lifecycle/' + $row.find('div.uuid').text(),
        method: 'GET',
        async: true,
        success: (result, status, xhr) => {
          // console.log('triggering with result', result)
          $form.trigger('rx', result)
        },
        error: (xhr, status, err) => {
          console.log(xhr, status, err)
        },
      })
      $row
        .parent()
        .find('.row.selected')
        .removeClass('selected')
      $row.addClass('selected')
    })

  $form.find('ul>li.attr.expandable').on('click', '>.label', e => {
    console.log('clicked', $(e.currentTarget)
      .parent()
      .toggleClass('expanded'))
  })

  $lifecycle
    .on('init', e => { })
    .on('activate', e => {
      $lifecycle
        .addClass('active')
        .find('>.ndx>.rows')
        .trigger('refresh')
    })
})