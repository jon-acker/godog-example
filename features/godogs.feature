Feature:

  As Library Member

  Scenario: Successfully borrowing a book
    Given "Jon" has registered as a member of "Hackney Library"
    When "Jon" tries to borrow the book "Harry Potter"
    Then the book "Harry Potter" should have been loaned to "Jon"

  Scenario: Failing to borrow a book when not a member
    Given "Jon" has not registered as a member of "Hackney Library"
    When "Jon" tries to borrow the book "Harry Potter"
    Then "Jon" should be told "Only members can borrow books"
