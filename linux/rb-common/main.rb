module TestUtils
  def HelloMessage
    "Hello, from Docker!"
  end
end


class RubyTest
  include TestUtils
end

rt = RubyTest.new

puts rt.HelloMessage
