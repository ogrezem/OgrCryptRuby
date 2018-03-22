# charset: UTF-8
=begin
Задача: везде всё сменить с двоичной системы на шеснадцатиричную
=end
require 'base64'
require 'digest'
require 'openssl'

module AESCrypt

  def AESCrypt.encrypt(password, iv, cleardata)
    cipher = OpenSSL::Cipher.new('AES-256-CBC')
    cipher.encrypt  # set cipher to be encryption mode
    cipher.key = password
    cipher.iv  = iv

    encrypted = ''
    encrypted << cipher.update(cleardata)
    encrypted << cipher.final
    AESCrypt.b64enc(encrypted)
  end

  def AESCrypt.decrypt(password, iv, secretdata)
    secretdata = Base64::decode64(secretdata)
    decipher = OpenSSL::Cipher::Cipher.new('aes-256-cbc')
    decipher.decrypt
    decipher.key = password
    decipher.iv = iv if iv != nil
    decipher.update(secretdata) + decipher.final
  end

  def AESCrypt.b64enc(data)
    Base64.encode64(data).gsub(/\n/, '')
  end

end

=begin
password = Digest::SHA256.digest('Nixnogen')
iv       = 'a2xhcgAAAAAAAAAA'
buf      = "Here is some data for the encrypt" # 32 chars
enc      = AESCrypt.encrypt(password, iv, buf)
dec      = AESCrypt.decrypt(password, iv, enc)
puts "encrypt length: #{enc.length}"
puts "encrypt in Base64: " + enc
puts "decrypt all: " + dec
=end

module Vernam
  def self.to_two arr
    result = []
    if arr.class == Array
      arr.each do |e|
        result << e.to_s(2).to_i
      end
    elsif arr.class == Fixnum
      result = arr.to_s(2).to_i
    end
    return result
  end

  def self.to_ten arr
    result = []
    arr.each do |e|
      result << e.to_s.to_i(2)
    end
    return result
  end

  def self.make_key len
    result = []
    len.times do |i|
      result << rand(1..255)
    end
    return result
  end

  def self.make_abc
    abc = []
    abc[0] = ([]<<0).pack('c*')
    255.times { |i| abc[i+1] = ([]<<(i+1)).pack('c*') }
    return abc
  end

  def self.crypt(v, k, sym)
    nb = []
    v.size.times do |b|
      if v[b] == 0
        nb << k[b]
      else
        bv = to_two(v[b])
        bk = to_two(k[b])
        nb << to_two("0b#{bv}".to_i(2) ^ "0b#{bk}".to_i(2))
      end
    end
    return nb if sym == :two
    return to_ten(nb) if sym == :ten
  end

  def self.crypt_v(v, k, arr)
    result = []
    ###################################################
    fgf = arr##
    #puts fgf.inspect##
    if fgf.include?(nil)
      puts 'YES'
    else
      puts 'NO'
    end
    ###################################################
    v.size.times do |c|
      begin
        rr = ((v[c] + k[c]) % 255) - 1
      rescue TypeError => error
        puts "+#{v[c]}+", "*#{k[c-2]}*", "*#{k[c-1]}*", "*#{k[c]}*", "*#{k[c+1]}*", "*#{k[c+2]}*", ''
        puts k.size, v.size, c
      end
      result << arr[rr].bytes[0]
=begin
      if (v[c] + k[c]) == 256
        result << arr[-1].bytes[0]
      elsif (v[c] + k[c]) == 255
        result << arr[-2].bytes[0]
      else
        result << arr[rr].bytes[0]
      end
=end
    end
    return result
  end

  def self.decrypt_v(v, k, arr)
    result = []
    v.size.times do |c|
      begin
        rr = (v[c] - k[c] + 255) % 255 + 1
      rescue TypeError => error
        puts v.size, k.size, c, "**#{k[c+2]}**", ''
      end
      begin
        result << arr[rr].bytes[0]
      rescue NoMethodError => error
        puts error.backtrace_locations
        puts rr
      end
    end
    return result
  end

  def self.crypt_b v
    new_v = []
    v << ' '.bytes[0] if (v.size % 2) != 0
    v.size.times do |c|
      if (c % 2) != 0
        new_v << v[c-1]
      elsif (c == 0) || ((c % 2) == 0)
        new_v << v[c+1]
      end
    end
    return new_v
  end

  def self.ogr_crypt(v, k, arr)
    result = crypt(crypt_v(crypt_b(v), k, arr), k, :ten)
  end

  def self.ogr_decrypt v, k, arr
    result = crypt_b(decrypt_v(crypt(v, k, :ten), k, arr))
  end

end